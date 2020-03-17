package net

import (
	"fmt"
  "net"
  "bufio"
  "io"
)

func dial() {
  conn, err := net.Dial("tcp", "127.0.0.1:8787")
  if err != nil {
    fmt.Print("error: ", err)
  }

  fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
  status, err := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("status: ", status)
}

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

func handleConnection(c net.Conn) {
  defer c.Close()
  io.Copy(c, c)
  fmt.Printf("Connection from %v closed.\n", c.RemoteAddr())
}