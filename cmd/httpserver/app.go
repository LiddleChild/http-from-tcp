package main

import (
	"fmt"
	"httpfromtcp/internal/http"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"strconv"
)

var (
	// id -> name
	users   = map[int]string{1000: "admin"}
	counter = 0
)

func ListUsers(w io.Writer, req *request.Request) *http.HandlerError {
	for id, name := range users {
		fmt.Fprintf(w, "%d: %s\n", id, name)
	}

	return nil
}

func UserById(w io.Writer, req *request.Request) *http.HandlerError {
	id, err := strconv.Atoi(req.Param["id"])
	if err != nil {
		return &http.HandlerError{
			Code:    response.StatusBadRequest,
			Message: "invalid id",
		}
	}

	name, ok := users[id]
	if !ok {
		return &http.HandlerError{
			Code:    response.StatusNotFound,
			Message: "user not found",
		}
	}

	_, err = w.Write([]byte(name))
	if err != nil {

	}

	return nil
}

func CreateUser(w io.Writer, req *request.Request) *http.HandlerError {
	username := string(req.Body)

	if len(username) == 0 {
		return &http.HandlerError{
			Code:    response.StatusBadRequest,
			Message: "username cannot be empty",
		}
	}

	users[counter] = username
	w.Write([]byte(strconv.Itoa(counter)))
	counter += 1

	return nil
}

func UpdateUser(w io.Writer, req *request.Request) *http.HandlerError {
	id, err := strconv.Atoi(req.Param["id"])
	if err != nil {
		return &http.HandlerError{
			Code:    response.StatusBadRequest,
			Message: "invalid id",
		}
	}

	oldName, ok := users[id]
	if !ok {
		return &http.HandlerError{
			Code:    response.StatusNotFound,
			Message: "user not found",
		}
	}

	name := string(req.Body)
	users[id] = name

	_, err = fmt.Fprintf(w, "%s updated to %s", oldName, name)
	if err != nil {

	}

	return nil

}

func DeleteUser(w io.Writer, req *request.Request) *http.HandlerError {
	id, err := strconv.Atoi(req.Param["id"])
	if err != nil {
		return &http.HandlerError{
			Code:    response.StatusBadRequest,
			Message: "invalid id",
		}
	}

	name, ok := users[id]
	if !ok {
		return &http.HandlerError{
			Code:    response.StatusNotFound,
			Message: "user not found",
		}
	}

	delete(users, id)

	_, err = fmt.Fprintf(w, "%s Deleted", name)
	if err != nil {

	}

	return nil
}
