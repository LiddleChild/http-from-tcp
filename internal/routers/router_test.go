package routers

import (
	"github.com/LiddleChild/http-from-tcp/internal/http"
	"github.com/LiddleChild/http-from-tcp/internal/request"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	r := NewRouter()

	defaultHandler := func(w io.Writer, req *request.Request) *http.HandlerError {
		return nil
	}

	var err error

	assert.Nil(t, r.GetHandler(http.MethodGet, "/"))
	err = r.RegisterHandler(http.MethodGet, "/", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/"))

	assert.Nil(t, r.GetHandler(http.MethodGet, "/a"))
	err = r.RegisterHandler(http.MethodGet, "/a", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a"))

	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a"))
	err = r.RegisterHandler(http.MethodGet, "/a", defaultHandler)
	assert.NotNil(t, err)

	assert.Nil(t, r.GetHandler(http.MethodGet, "/b"))
	err = r.RegisterHandler(http.MethodGet, "/b", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/b"))

	assert.Nil(t, r.GetHandler(http.MethodGet, "/a/b"))
	err = r.RegisterHandler(http.MethodGet, "/a/b", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a/b"))

	assert.Nil(t, r.GetHandler(http.MethodGet, "/a/b/c/d/e"))
	err = r.RegisterHandler(http.MethodGet, "/a/b/c/d/e", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a/b/c/d/e"))

	assert.Nil(t, r.GetHandler(http.MethodGet, "/a/:b/c"))
	err = r.RegisterHandler(http.MethodGet, "/a/:b/c", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a/:b/c"))

	err = r.RegisterHandler(http.MethodGet, "/a/:c/c", defaultHandler)
	assert.NotNil(t, err)

	assert.Nil(t, r.GetHandler(http.MethodGet, "/a/:b/d"))
	err = r.RegisterHandler(http.MethodGet, "/a/:b/d", defaultHandler)
	assert.Nil(t, err)
	assert.NotNil(t, r.GetHandler(http.MethodGet, "/a/:b/d"))
}
