package main

import (
    "fmt"
    "net"
    "runtime"
    "strconv"
    "time"
)

type  portpair struct {
    first uint16
    last uint16
}

func main() {
    exclusions := []portpair{{49692,49791 }, {49792, 49891}, {49892, 49991}, {49992, 50091}, {50092, 50191}, {50214, 50313}, {50498, 50597}, {}}
    fmt.Printf("%d\n", len(exclusions))
    fmt.Printf("OS? %v \n", runtime.GOOS)

    /*ln, _ := net.Listen("tcp", "localhost:0")
    fmt.Printf("%s %s \n", ln.Addr().Network(), ln.Addr().String())
    ln.Close()

    ln, _ = net.Listen("tcp4", "localhost:0")
    fmt.Printf("%s %s \n", ln.Addr().Network(), ln.Addr().String())
    ln.Close()

    ln, _ = net.Listen("tcp6", "localhost:0")
    fmt.Printf("%s %s \n", ln.Addr().Network(), ln.Addr().String())
    ln.Close()*/

    excluded := 0
    isFirstPort := true
    for i:=0; i < 100; i++ {
        p := GetAvailablePort()
        fmt.Printf("Got port %d\n", p)
        if isFirstPort {
            isFirstPort = false
            newExclude := portpair{p, p+15}
            exclusions = append(exclusions, newExclude)
        }

        for _, pair := range exclusions {
            if p >= pair.first && p <= pair.last {
                fmt.Printf(">>>>>>>>> Excluded %d beause of range %d to %d\n", p, pair.first, pair.last)
                excluded++
            }
        }

        time.Sleep(100 * time.Millisecond)
    }


    fmt.Printf("Excluded %d entries\n", excluded)
}

func GetAvailablePort() uint16 {
    endpoint := GetAvailableLocalAddress()
    _, port, err := net.SplitHostPort(endpoint)
    die(err)

    portInt, err := strconv.Atoi(port)
    die(err)

    return uint16(portInt)
}

func GetAvailableLocalAddress() string {
    ln, err := net.Listen("tcp", "localhost:0")
    die(err)
    // There is a possible race if something else takes this same port before
    // the test uses it, however, that is unlikely in practice.
    defer ln.Close()
    return ln.Addr().String()
}

func die(e error) {
    if e != nil {
        panic(e)
    }
}
