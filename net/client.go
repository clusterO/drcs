package network

import (
	"bufio"
	"fmt"
	"net"
)

func GetVersions(result string) {
	// remote callback
}

func GotVersions(result string) {
	print("server: ", result)
}

func Dial(ip string, port string) {
	conn, err := net.Dial("tcp", ip + ":" + port)
	if err != nil {
		fmt.Print("error: ", err)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("status: ", status)
}

func Connect(ip string, port string, direcotry string, op bool) {

}
