package naspad

import (
	"net/http"
)

type Server struct {
	router *Router
}

func NewServer() *Server {
	return &Server{router: NewRouter()}
}

func (s *Server) Handle(method, path string, handler HandlerFunc) {
	s.router.Handle(method, path, handler)
}

func (s *Server) Use(middleware MiddlewareFunc) {
	s.router.Use(middleware)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
