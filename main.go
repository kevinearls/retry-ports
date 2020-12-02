package main

import (
    "fmt"
    "runtime"
)

func main() {
    //exclusions := []portpair{{"49692","49791" }, {"49792", "49891"}, {"49892", "49991"}, {"49992", "50091"}, {"50092", "50191"}, {"50214", "50313"}, {"50498", "50597"}, {}}
    //fmt.Printf("%d\n", len(exclusions))
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

