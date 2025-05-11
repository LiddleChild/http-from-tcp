package main

import (
	"fmt"
	"github.com/LiddleChild/http-from-tcp/internal/request"
	"net"
)

const host = "0.0.0.0:42069"

func main() {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Listening TCP on %v\n", host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			panic(err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	}
}
