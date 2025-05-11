package server

import (
	"bytes"
	"fmt"
	"github.com/LiddleChild/http-from-tcp/internal/http"
	"github.com/LiddleChild/http-from-tcp/internal/request"
	"github.com/LiddleChild/http-from-tcp/internal/response"
	"github.com/LiddleChild/http-from-tcp/internal/routers"
	"github.com/LiddleChild/http-from-tcp/internal/worker"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type Server struct {
	listener net.Listener
	router   *routers.Router
	pool     *worker.Pool[net.Conn]

	closed *atomic.Bool
	err    error
}

func NewServer() *Server {
	b := &atomic.Bool{}
	b.Store(false)

	server := &Server{
		listener: nil,
		router:   routers.NewRouter(),
		pool:     worker.NewPool[net.Conn](1024),
		closed:   b,
		err:      nil,
	}

	server.pool.TaskFunc(server.handle)

	return server
}

func (s *Server) Serve(address string) error {
	s.router.ListRoutes()

	if s.err != nil {
		return fmt.Errorf("server error: %w", s.err)
	}

	var err error

	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		s.pool.QueueTask(conn)
	}

	return nil
}

func (s *Server) handle(conn net.Conn) {
	start := time.Now()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var (
		bs   []byte
		code response.StatusCode

		buffers = bytes.NewBuffer([]byte{})
		handler = s.router.GetHandler(http.Method(req.RequestLine.Method), req.RequestLine.RequestTarget)
	)

	writeResponse := func() {
		if err := response.WriteStatusLine(conn, code); err != nil {
			fmt.Println(err)
			return
		}

		headers := response.GetDefaultHeader(len(bs))
		if err := response.WriteHeaders(conn, headers); err != nil {
			fmt.Println(err)
			return
		}

		_, err = conn.Write(bs)
		if err != nil {
			fmt.Println(err)
			return
		}

		duration := time.Since(start).Round(time.Microsecond)
		log.Printf("%d | %12s | %-6s %s\n", code, duration.String(), req.RequestLine.Method, req.RequestLine.RequestTarget)
	}

	s.router.ParseParam(req.Param, req.RequestLine.RequestTarget)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic recovered: %v\n", r)

			bs = []byte("internal server error")
			code = response.StatusInternalServerError

			writeResponse()
		}
	}()

	if handler == nil {
		bs = []byte("404 not found")
		code = response.StatusNotFound
	} else if resp := handler(buffers, req); resp != nil {
		bs = []byte(resp.Message)
		code = resp.Code
	} else {
		bs = buffers.Bytes()
		code = response.StatusOK
	}

	writeResponse()
}

func (s *Server) Close() error {
	s.closed.Store(true)

	s.pool.Close()

	if err := s.listener.Close(); err != nil {
		return err
	}

	return nil
}
