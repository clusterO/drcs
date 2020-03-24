package net

import (
	"bufio"
	"fmt"
	"net"
)

func getVersions(result string) {
	// remote callback
}

func gotVersions(result string) {
	print("server: ", result)
}

func dial() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Print("error: ", err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("status: ", status)
}
