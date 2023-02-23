package main

import (
	"fmt"
	"net"
)

// client test
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Print("client init failed.")
		return
	}
	conn.Write([]byte("Halo, World."))
}
