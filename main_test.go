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

func TestOne(t *testing.T) {
	fmt.Printf("Runtime.GOOS is %v\n", runtime.GOOS)
	for i:=0; i < 25 ; i++ {
		ln, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			t.Errorf("net.Listen got error %v\n", err)
		} else {
			fmt.Printf(">>>> Got address %s\n", ln.Addr().String())
			ln.Close()
		}
	}
	fmt.Println("It worked!!!")

	if runtime.GOOS == "windows" {
		fmt.Println(">>>>> We're on windows")
		cmd := exec.Command(/*"cmd", "/C",*/ "netsh", "interface",  "ipv4",  "show",  "excludedportrange", "protocol=tcp")
		output, err := cmd.CombinedOutput()
		if err != nil  {
			fmt.Errorf("netsh command got error %v\n", err)
		}

		fmt.Printf("NETSH command got: \n%s\n", string(output))
	}
}

func TestExclusions(t *testing.T) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(/*"cmd", "/C",*/ "netsh", "interface",  "ipv4",  "show",  "excludedportrange", "protocol=tcp")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		exclusions := getExclusionsList(string(output), t)
		fmt.Printf(">>>>>> Got %d exclusion pairs\n", len(exclusions))
		for _, p := range exclusions {
			fmt.Printf("\t%v\n", p)
		}

		fmt.Printf("NETSH command got: \n%s\n", string(output))
	}
}

func TestTwo(t *testing.T) {
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
			//fmt.Printf("Pair: %s, %s\n", entries[0], entries[1])
			first, err := strconv.Atoi(entries[0])
			require.NoError(t, err)
			second, _ := strconv.Atoi(entries[1])
			require.NoError(t, err)
			pair := portpair{first: uint16(first), last:  uint16(second),}
			exclusions = append(exclusions, pair)
		}
	}
	return exclusions
}

