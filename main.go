package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server is the main entity that will handle all request. It implements
// "http.Handler" and "capo.Group" interfaces.
type Server struct {
	engine *gin.Engine
}

// New creates a new server instance.
func New() *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	return &Server{
		engine: engine,
	}
}

// ServeHTTP implements "http.Handler" interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}

// Use adds a handler function in the request handlers chain.
func (s *Server) Use(handlers ...http.HandlerFunc) {
	h := handlersToGinHandlers(handlers...)
	s.engine.Use(h...)
}

// Group creates a new group to handle http requests.
func (s *Server) Group(relativePath string, handlers ...http.HandlerFunc) Group {
	h := handlersToGinHandlers(handlers...)
	ginG := s.engine.Group(relativePath, h...)
	return newGroup(ginG)
}

// Get handles a GET request.
func (s *Server) Get(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	s.engine.GET(relativePath, h)
}

// Get handles a POST request.
func (s *Server) Post(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	s.engine.POST(relativePath, h)
}

// Get handles a PUT request.
func (s *Server) Put(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	s.engine.PUT(relativePath, h)
}

// Get handles a DELETE request.
func (s *Server) Delete(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	s.engine.DELETE(relativePath, h)
}

func handlersToGinHandlers(handlers ...http.HandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		ginHandlers[i] = gin.WrapF(h)
	}

	return ginHandlers
}
