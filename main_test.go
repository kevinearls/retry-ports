package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
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
		cmd := exec.Command("cmd", "/C", "netsh", "interface", "ipv4", "show", "excludedportrange protocol=tcp")
		output, err := cmd.CombinedOutput()
		if err != nil  {
			fmt.Errorf("netsh command got error %v\n", err)
		}

		fmt.Printf("NETSH command got: \n%s\n", string(output))
	}
}
