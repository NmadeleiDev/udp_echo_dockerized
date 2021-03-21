package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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

type testMsg struct {
	PeerPort	uint16	`json:"peer_port"`
	Data		string	`json:"data"`
}

func makeTestTcpCall(addr string) {
	fmt.Printf("Calling %v", addr)

	client := http.Client{Timeout: time.Second * 5}

	resp, err := client.Get(addr)
	if err != nil {
		fmt.Printf("Error sending get peer req: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Resp status invalid: %v", resp.StatusCode)
	} else {
		fmt.Printf("Peer get call success! %v %v", resp.StatusCode, resp.Status)
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

		var msg testMsg

		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			fmt.Printf("Error unmarshal: %v", err)
			return
		}

		fmt.Printf("%d bytes received: '%s' from: %s\n%s\nPort: %v; msg: %v;\n", n, string(buf[:n]), addr.String(), time.Now(),
			msg.PeerPort, msg.Data)

		makeTestTcpCall(fmt.Sprintf("http://%s:%d/", addr.IP.String(), msg.PeerPort))

		fmt.Printf("Writing response to %v\n", *addr)

		if err != nil {
			fmt.Println("error: ", err)
		}

		_, _ = ServerConn.WriteTo([]byte(msg.Data), addr)
	}
}