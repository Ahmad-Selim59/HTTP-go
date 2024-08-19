package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type HttpResponse struct {
	Status  string
	Headers map[string]string
	Body    string
}

func main() {
	response, error := sendHttpRequest("34.76.19.201", 8080, "/sessions", "GET", map[string]string{}, "")
	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println(response)
	}
}

func sendHttpRequest(address string, portNum int, path string, methodType string, headers map[string]string, requestBody string) (*HttpResponse, error) {
	url := fmt.Sprintf("%s:%d", address, portNum)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	requestLine := fmt.Sprintf("%s %s HTTP/1.0\r\n", methodType, path)
	if methodType == "POST" || methodType == "PUT" {
		headers["Content-Length"] = fmt.Sprintf("%d", len(requestBody))
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/x-www-form-urlencoded" // Default type if not set
		}
	}

	fullRequest := requestLine
	for key, value := range headers {
		fullRequest += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	fullRequest += "\r\n" // End of headers
	if methodType == "POST" || methodType == "PUT" {
		fullRequest += requestBody
	}

	fmt.Fprint(conn, fullRequest)

	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	headersMap := make(map[string]string)

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
			headersMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	var body string
	if contentLength, ok := headersMap["Content-Length"]; ok {
		fmt.Println("Content-Length:", contentLength)
		contentLengthInt, _ := strconv.Atoi(contentLength)
		bodyBytes := make([]byte, contentLengthInt)
		_, err := io.ReadFull(reader, bodyBytes)
		if err != nil && err != io.EOF {
			return nil, err
		}
		body = string(bodyBytes)
	} else {
		// Fallback if Content-Length doesn't work or isn't available
		var bodyBuilder strings.Builder
		for {
			part := make([]byte, 1024)
			n, err := reader.Read(part)
			if n > 0 {
				bodyBuilder.Write(part[:n])
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
		}
		body = bodyBuilder.String()
	}

	response := &HttpResponse{
		Status:  strings.TrimSpace(status),
		Headers: headersMap,
		Body:    body,
	}

	return response, nil
}
