/**
 * Based on https://github.com/docker/for-win/issues/3171
 *
 * On Windows it appears certain ports get reserved by Hyper-V and net.Listen("tcp", "localhost:0") may return a port
 * within one of those ranges
 */

package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestExclusions(t *testing.T) {
	if runtime.GOOS == "windows" {
		//exclusions := createExclustionsList(t)
		//fmt.Printf(">>>>>> Got %d exclusion pairs\n", len(exclusions))
		//for _, p := range exclusions {
		//	fmt.Printf("\t%v\n", p)
		//}

		port := GetAvailablePort(t)
		fmt.Printf("Got first port %d\n", port)

		// HACK!  Add something we know will cause exclusions
		//newExclude := portpair{strconv.Itoa(int(port)), strconv.Itoa(int(port)+15)}
		//exclusions = append(exclusions, newExclude)
		//fmt.Printf("Added %v to exclusions\n", newExclude)

		secondPort := GetAvailablePort(t)
		fmt.Printf("Got seondPort %d\n", secondPort)
		require.Equal(t, port + 16, secondPort)
	} else {
		port := GetAvailablePort(t)
		fmt.Printf("Got first port %d\n", port)
		secondPort := GetAvailablePort(t)
		fmt.Printf("Got seondPort %d\n", secondPort)
		require.Equal(t, port + 1, secondPort)  // Is this always true?
	}
}

type  portpair struct {
	first string
	last string
}

func TestHardcodedExclusions(t *testing.T) {
	// If emtpy it looks like this:
	/*

	Protocol tcp Port Exclusion Ranges

	Start Port    End Port
	----------    --------

	* - Administered port exclusions.
	 */
	emptyExclusionsText :=`

Protocol tcp Port Exclusion Ranges

Start Port    End Port      
----------    --------      

* - Administered port exclusions.`

	exclusionsText := `

Start Port    End Port
----------    --------
     49697       49796
     49797       49896

* - Administered port exclusions.
`
	exclusions := getExclusionsList(exclusionsText, t)
	fmt.Printf("Added %d pairs to exclusion list\n", len(exclusions))
	for _, p := range exclusions {
		fmt.Printf("\t%v\n", p)
	}

	emptyExclusions := getExclusionsList(emptyExclusionsText, t)
	fmt.Printf("Empty got %d pairs\n", len(emptyExclusions))
}

func createExclustionsList(t *testing.T) []portpair {
	cmd := exec.Command("netsh", "interface",  "ipv4",  "show",  "excludedportrange", "protocol=tcp")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	fmt.Printf("NETSH command got: \n%s\n", string(output))   //FIXME remove

	exclusions := getExclusionsList(string(output), t)

	/// FIXME FIXME FIXME remove, hack for testing
	// HACK!  Add something we know will cause exclusions
	//
	var port string
	endpoint := GetAvailableLocalAddress(t)
	_, port, err = net.SplitHostPort(endpoint)
	require.NoError(t, err)

	stupid, _ := strconv.Atoi(port)
	p2:= strconv.Itoa(stupid + 10)
	newExclude := portpair{port, p2}
	exclusions = append(exclusions, newExclude)
	fmt.Printf("Added %v to exclusions\n", newExclude)
	// End of HACK

	return exclusions
}

// Get excluded ports on Windows from the command: netsh interface ipv4 show excludedportrange protocol=tcp
func getExclusionsList(exclusionsText string, t *testing.T) []portpair {
	exclusions := []portpair{}

	parts := strings.Split(exclusionsText, "--------")
	require.Equal(t, len(parts), 3)
	portsText := strings.Split(parts[2], "*")
	require.Equal(t, len(portsText), 2)
	lines := strings.Split(portsText[0], "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			entries := strings.Fields(strings.TrimSpace(line))
			require.Equal(t, len(entries), 2)
			pair := portpair{entries[0], entries[1]}
			exclusions = append(exclusions, pair)
		}
	}
	return exclusions
}


func GetAvailablePort(t *testing.T) uint16 {
	var exclusions [] portpair
	if runtime.GOOS == "windows" {
		exclusions = createExclustionsList(t)
	}

	portFound := false
	var port string
	var err error
	for !portFound {
		endpoint := GetAvailableLocalAddress(t)
		_, port, err = net.SplitHostPort(endpoint)
		require.NoError(t, err)
		portFound = true
		if runtime.GOOS == "windows" {
			for _, pair := range exclusions {
				if port >= pair.first && port <= pair.last {
					fmt.Printf(">>>>>>>>> Excluded %s because of range %s to %s\n", port, pair.first, pair.last)  // TODO change to debug line or remove
					portFound = false
					break
				}
			}
		}
	}

	portInt, err := strconv.Atoi(port)
	require.NoError(t, err)

	return uint16(portInt)
}

func GetAvailableLocalAddress(t *testing.T) string {
	ln, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	// There is a possible race if something else takes this same port before
	// the test uses it, however, that is unlikely in practice.
	defer ln.Close()
	return ln.Addr().String()
}

