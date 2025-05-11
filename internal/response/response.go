package response

import (
	"fmt"
	"github.com/LiddleChild/http-from-tcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	status := ""
	switch statusCode {
	case StatusOK:
		status = "HTTP/1.1 200 OK\n"
	case StatusBadRequest:
		status = "HTTP/1.1 400 Bad Request\n"
	case StatusNotFound:
		status = "HTTP/1.1 404 Not Found\n"
	case StatusInternalServerError:
		status = "HTTP/1.1 500 Internal Server Error\n"
	}

	_, err := w.Write([]byte(status))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeader(contentLength int) headers.Headers {
	header := headers.NewHeaders()

	header["Content-Length"] = strconv.Itoa(contentLength)
	header["Connection"] = "close"
	header["Content-Type"] = "text/plain"

	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s\n", key, value)
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}
