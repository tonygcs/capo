package capo

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler is the data type for the method that will handle the HTTP requests.
type Handler func(*Context) error

// Server is the main entity that will handle all request. It implements
// "http.Handler" and "capo.Group" interfaces.
type Server struct {
	g Group
	r *mux.Router
}

// New creates a new server instance.
func New() *Server {
	router := mux.NewRouter()
	return &Server{
		g: newGroup("/", nil, router),
		r: router,
	}
}

// Default returns a server with the default configuration.
func Default() *Server {
	s := New()
	s.UseAfterAlways(ErrorHandling)
	return s
}

// ServeHTTP implements "http.Handler" interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

// Use adds the handlers that will run before each request. Note that if a
// handler returns an error, the next handlers won't run.
func (s *Server) UseBefore(handlers ...Handler) {
	s.g.UseBefore(handlers...)
}

// UseAfter adds the handlers that will run after each request if it is not cancelled.
func (s *Server) UseAfter(handlers ...Handler) {
	s.g.UseAfter(handlers...)
}

// UseAfterAlways adds the handlers that will run after every request.
func (s *Server) UseAfterAlways(handlers ...func(*Context)) {
	s.g.UseAfterAlways(handlers...)
}

// Group creates a new group to handle http requests.
func (s *Server) Group(relativePath string) Group {
	return s.g.Group(relativePath)
}

// Get handles a GET request.
func (s *Server) Get(relativePath string, handler Handler) {
	s.g.Get(relativePath, handler)
}

// Get handles a POST request.
func (s *Server) Post(relativePath string, handler Handler) {
	s.g.Post(relativePath, handler)
}

// Get handles a PUT request.
func (s *Server) Put(relativePath string, handler Handler) {
	s.g.Put(relativePath, handler)
}

// Get handles a DELETE request.
func (s *Server) Delete(relativePath string, handler Handler) {
	s.g.Delete(relativePath, handler)
}

func (s *Server) path() string {
	return s.g.path()
}

func (s *Server) getBeforeHandlers() []Handler {
	return s.g.getBeforeHandlers()
}

func (s *Server) getAfterHandlers() []Handler {
	return s.g.getAfterHandlers()
}

func (s *Server) getAfterAlwaysHandlers() []func(*Context) {
	return s.g.getAfterAlwaysHandlers()
}
