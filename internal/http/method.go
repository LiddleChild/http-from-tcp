package http

import (
	"github.com/LiddleChild/http-from-tcp/internal/request"
	"github.com/LiddleChild/http-from-tcp/internal/response"
	"io"
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPost   Method = "POST"
	MethodPut    Method = "PUT"
	MethodPatch  Method = "PATCH"
	MethodDelete Method = "DELETE"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
