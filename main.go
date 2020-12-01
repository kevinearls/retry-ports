package main

import (
    "fmt"
    "github.com/stretchr/testify/require"
    "net"
    "runtime"
    "strconv"
    "testing"
)

type  portpair struct {
    first string
    last string
}

func main() {
    exclusions := []portpair{{"49692","49791" }, {"49792", "49891"}, {"49892", "49991"}, {"49992", "50091"}, {"50092", "50191"}, {"50214", "50313"}, {"50498", "50597"}, {}}
    fmt.Printf("%d\n", len(exclusions))
    fmt.Printf("OS? %v \n", runtime.GOOS)

    /*excluded := 0
    isFirstPort := true

    for i:=0; i < 30; i++ {
        p := GetAvailablePort(exclusions)
        fmt.Printf("Got port %d\n", p)
        if isFirstPort {  // Hack add more excluded ports to test condition
            isFirstPort = false
            newExclude := portpair{strconv.Itoa(int(p)), strconv.Itoa(int(p)+15)}
            exclusions = append(exclusions, newExclude)
            fmt.Printf("Added %v to exclusions\n", newExclude)
        }
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Printf("Excluded %d entries\n", excluded) */
}

// This should probably get the exclusions list itself if it's on windows?????
func GetAvailablePort(t *testing.T, exclusions []portpair) uint16 {
    portFound := false
    var port string
    var err error
    for !portFound {
        endpoint := GetAvailableLocalAddress(t)
        _, port, err = net.SplitHostPort(endpoint)
        require.NoError(t, err)
        portFound = true
        if runtime.GOOS == "windows" {   // FIXME FIXME change back to is windows!
            for _, pair := range exclusions {
                if port >= pair.first && port <= pair.last {
                    fmt.Printf(">>>>>>>>> Excluded %s because of range %s to %s\n", port, pair.first, pair.last)  // TODO change to debug line
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
