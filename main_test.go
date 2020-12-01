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
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestExclusions(t *testing.T) {
	if runtime.GOOS == "windows" {
		exclusions := createExclustionsList(t)
		fmt.Printf(">>>>>> Got %d exclusion pairs\n", len(exclusions))
		for _, p := range exclusions {
			fmt.Printf("\t%v\n", p)
		}

		port := GetAvailablePort(t, exclusions)
		fmt.Printf("Got first port %d\n", port)

		// HAK!  Add something we know will cause exclusions
		newExclude := portpair{strconv.Itoa(int(port)), strconv.Itoa(int(port)+15)}
		exclusions = append(exclusions, newExclude)
		fmt.Printf("Added %v to exclusions\n", newExclude)

		secondPort := GetAvailablePort(t, exclusions)
		fmt.Printf("Got seondPort %d\n", secondPort)
		require.Equal(t, port + 16, secondPort)
	}
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
	fmt.Printf("NETSH command got: \n%s\n", string(output))

	exclusions := getExclusionsList(string(output), t)
	return exclusions
}

// Get excluded ports on Windows from the command: netsh interface ipv4 show excludedportrange protocol=tcp
func getExclusionsList(exclusionsText string, t *testing.T) []portpair {
	exclusions := []portpair{}

	parts := strings.Split(exclusionsText, "--------")
	require.Equal(t, len(parts), 3)
	portsText := strings.Split(parts[2], "*")  // TODO check for two
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

