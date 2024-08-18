package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	tcpServer("httpbin.org", 80)
}

func tcpServer(address string, portNum int) {
	url := fmt.Sprintf("%s:%d", address, portNum)
	conn, err := net.Dial("tcp", url)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "GET /anything HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(status)
	}
}
