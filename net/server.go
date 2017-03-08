package net

import (
	"fmt"
	"net"
	"bufio"
	"io"
	"dcrs/dcrs"
)

func listen() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Print("error: ", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
		fmt.Print("error: ", err)
		}
		
		go handleConnection(conn)
	}
}

func remoteGetCommits(directory string) {
  	return GetAllCommits()
}

func remoteGetCommitsContent(directory string, commit string) {
  	return CompressAndSend()
}

func handleConnection(c net.Conn) {
	defer c.Close()
	io.Copy(c, c)
	fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}