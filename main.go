package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	tcpServer("httpbin.org", 80, "/anything", "GET")
}

func tcpServer(address string, portNum int, path string, methodType string) {
	url := fmt.Sprintf("%s:%d", address, portNum)
	conn, err := net.Dial("tcp", url)

	if err != nil {
		log.Fatalf("failed to connect to %s on port %d. Received the following error %v", address, portNum, err)
		return
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	fmt.Fprintf(conn, "%s %s HTTP/1.0\r\n\r\n", methodType, path)

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read status line - %v", err)
	}

	fmt.Printf("status: %s", status)

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error encountered %v", err)
			break
		}
		line = strings.TrimSpace(line)

		if line == "" {
			break //reached end
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	fmt.Println("Headers:")
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
	}

	if contentLength, ok := headers["Content-Length"]; ok {
		fmt.Printf("Body (first %s bytes):\n", contentLength)
		body := make([]byte, 1024)
		n, _ := reader.Read(body)
		fmt.Println(string(body[:n]))
	} else {
		fmt.Println("No Content-Length specified or no body")
	}
}
