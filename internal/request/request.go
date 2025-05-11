package request

import (
	"bytes"
	"errors"
	"github.com/LiddleChild/http-from-tcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const (
	BufferSize = 8
)

type RequestState int

const (
	RequestStateInitialized RequestState = iota
	RequestStateParsingHeaders
	RequestStateParsingBody
	RequestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	Param       map[string]string

	state RequestState
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case RequestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		// early return for not enough data
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = RequestStateParsingHeaders

		return n, nil
	case RequestStateParsingHeaders:
		count := 0
		for {
			terminator := bytes.Index(data[count:], []byte("\r\n"))
			if terminator == 0 {
				count += 2
				break
			}

			n, done, err := r.Headers.Parse(data[count:])
			if err != nil {
				return 0, err
			}

			if !done {
				return count, nil
			}

			count += n
		}

		if _, ok := r.Headers["content-length"]; !ok {
			r.state = RequestStateDone
		} else {
			r.state = RequestStateParsingBody
		}

		return count, nil
	case RequestStateParsingBody:
		contentLength, err := strconv.Atoi(r.Headers["content-length"])
		if err != nil {
			return 0, errors.New("content-length is not an integer")
		}

		if contentLength == 0 {
			r.state = RequestStateDone
			return 0, nil
		}

		// early return since we dont know if the body has ended
		if len(data) < contentLength {
			return 0, nil
		} else if len(data) > contentLength {
			return 0, errors.New("body is bigger than content-length")
		}

		r.Body = make([]byte, len(data))
		copy(r.Body, data)

		r.state = RequestStateDone

		return contentLength, nil
	case RequestStateDone:
		return 0, errors.New("request is done parsing")
	default:
		return 0, errors.New("unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var (
		req = &Request{
			state:   RequestStateInitialized,
			Headers: headers.NewHeaders(),
			Body:    nil,
			Param:   make(map[string]string),
		}
		buffer     = make([]byte, BufferSize)
		startIndex = 0
	)

	for req.state != RequestStateDone {
		// if buffer is full (or greater just in case) double the buffers size
		if startIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		n, err := reader.Read(buffer[startIndex:])
		if errors.Is(err, io.EOF) {
			if req.state != RequestStateDone {
				return nil, errors.New("invalid request")
			}

			req.state = RequestStateDone
			break
		} else if err != nil {
			return nil, err
		}

		// accumulate how much bytes are read
		startIndex += n

		count := 0
		for req.state != RequestStateDone {
			n, err = req.parse(buffer[count:startIndex])
			if err != nil {
				return nil, err
			}

			if n == 0 {
				break
			}

			count += n
		}

		// shift all buffers to the left to prevent from keeps expanding the buffers
		copy(buffer, buffer[count:])
		startIndex -= count
	}

	return req, nil
}

func parseRequestLine(bs []byte) (*RequestLine, int, error) {
	terminator := bytes.Index(bs, []byte("\r\n"))
	if terminator == -1 {
		return nil, 0, nil
	}

	parts := strings.Split(string(bs[:terminator]), " ")
	if len(parts) < 3 {
		return nil, 0, errors.New("too few parts")
	}

	req := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}

	// method
	if req.Method != strings.ToUpper(req.Method) {
		return nil, 0, errors.New("method can only contains upper characters")
	}

	switch req.Method {
	case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
		break
	default:
		return nil, 0, errors.New("invalid method")
	}

	// request target
	if !strings.HasPrefix(req.RequestTarget, "/") {
		return nil, 0, errors.New("invalid request target")
	}

	// version
	if req.HttpVersion != "1.1" {
		return nil, 0, errors.New("unsupported protocol")
	}

	// terminator is an index of \r\n, so to return processed string it needs to include \r and \n characters
	return req, terminator + 2, nil
}
