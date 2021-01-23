package network

import (
	"fmt"
	"io"
	"net"
)

func listen() {
	ln, err := net.Listen("tcp", ":9999"); if err != nil {
		fmt.Print("error: ", err)
	}

	for {
		conn, err := ln.Accept(); if err != nil {
			fmt.Print("error: ", err)
		}
		go handleConnection(conn)
	}
}

func remoteGetCommits(directory string) string {
	return  ""
	// dcrs.GetAllCommits()
}

func remoteGetCommitsContent(directory string, commit string) string {
	return ""
	// dcrs.CompressAndSend("")
}

func handleConnection(c net.Conn) {
	defer c.Close()
	io.Copy(c, c)
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}
