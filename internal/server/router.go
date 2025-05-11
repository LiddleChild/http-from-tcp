package server

import "github.com/LiddleChild/http-from-tcp/internal/http"

func (s *Server) route(method http.Method, path string, handler http.Handler) {
	if s.err != nil {
		return
	}

	s.err = s.router.RegisterHandler(method, path, handler)
}

func (s *Server) Get(path string, handler http.Handler) {
	s.route(http.MethodGet, path, handler)
}

func (s *Server) Post(path string, handler http.Handler) {
	s.route(http.MethodPost, path, handler)
}

func (s *Server) Put(path string, handler http.Handler) {
	s.route(http.MethodPut, path, handler)
}

func (s *Server) Patch(path string, handler http.Handler) {
	s.route(http.MethodPatch, path, handler)
}

func (s *Server) Delete(path string, handler http.Handler) {
	s.route(http.MethodDelete, path, handler)
}
