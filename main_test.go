package main

import (
	"fmt"
	"net"
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
}
