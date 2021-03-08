package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// CheckError checks for errors
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	port := os.Getenv("LISTEN_PORT")
	ServerAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:" + port)
	CheckError(err)
	fmt.Println("listening on ", *ServerAddr)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		echoPort := binary.BigEndian.Uint32(buf[:4])
		addr.Port = int(echoPort)

		fmt.Printf("Writing response to %v", *addr)

		fmt.Printf("received: '%s' from: %s\n%v\n", string(buf[4:n]), addr.String(), time.Now())

		if err != nil {
			fmt.Println("error: ", err)
		}
		
		ServerConn.WriteTo(buf[0:n], addr)
	}
}