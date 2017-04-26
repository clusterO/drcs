package network

import (
	"fmt"
	"io"
	"net"
)

func Listen() {
	// should return when err
	ln, err := net.Listen("tcp4", "127.0.0.1:8181"); if err != nil {
		fmt.Print("error: ", err)
	}

	for 
	{
		conn, err := ln.Accept(); if err != nil {
			fmt.Print("error: ", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	io.Copy(c, c)
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}
