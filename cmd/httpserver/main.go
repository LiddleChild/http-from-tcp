package main

import (
	"httpfromtcp/internal/http"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const address = "0.0.0.0:42069"

func main() {
	server := server.NewServer()

	server.Get("/", func(w io.Writer, req *request.Request) *http.HandlerError { return nil })
	server.Get("/users", ListUsers)
	server.Post("/users", CreateUser)
	server.Get("/users/:id", UserById)
	server.Patch("/users/:id", UpdateUser)
	server.Delete("/users/:id", DeleteUser)

	go func() {
		defer server.Close()
		if err := server.Serve(address); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Println("Server started on port", address)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("Server gracefully stopped")
}
