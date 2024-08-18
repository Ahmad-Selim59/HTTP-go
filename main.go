package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type HttpResponse struct {
	Status  string
	Headers map[string]string
	Body    string
}

func main() {
	response, error := sendHttpRequest("httpbin.org", 80, "/anything", "GET")
	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println(response)
	}
}

func sendHttpRequest(address string, portNum int, path string, methodType string) (*HttpResponse, error) {
	url := fmt.Sprintf("%s:%d", address, portNum)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	fmt.Fprintf(conn, "%s %s HTTP/1.0\r\nHost: %s\r\n\r\n", methodType, path, address)

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)

		if line == "" {
			break // reached end of headers
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	var body string
	if _, ok := headers["Content-Length"]; ok {
		bodyBytes := make([]byte, 1024)
		n, _ := reader.Read(bodyBytes)
		body = string(bodyBytes[:n])
	}

	response := &HttpResponse{
		Status:  strings.TrimSpace(status),
		Headers: headers,
		Body:    body,
	}

	return response, nil
}
