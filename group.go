package capo

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"

	"github.com/gorilla/mux"
)

type Group interface {
	path() string
	getBeforeHandlers() []Handler
	getAfterHandlers() []Handler
	getAfterAlwaysHandlers() []func(*Context)

	// Use adds the handlers that will run before each request. Note that if a
	// handler returns an error, the next handlers won't run.
	UseBefore(handlers ...Handler)
	// UseAfter adds the handlers that will run after each request if it is not cancelled.
	UseAfter(handlers ...Handler)
	// UseAfterAlways adds the handlers that will run after every request.
	UseAfterAlways(handlers ...func(*Context))
	// Group creates a new group to handle http requests.
	Group(relativePath string) Group

	// Get handles a GET request.
	Get(relativePath string, handler Handler)
	// Get handles a POST request.
	Post(relativePath string, handler Handler)
	// Get handles a PUT request.
	Put(relativePath string, handler Handler)
	// Get handles a DELETE request.
	Delete(relativePath string, handler Handler)
}

// group is the group to wrap http handlers.
type group struct {
	mu sync.Mutex

	groupPath   string
	router      *mux.Router
	parent      Group
	children    []Group
	before      []Handler
	after       []Handler
	afterAlways []func(*Context)
}

// newGroup creates a new group instance.
func newGroup(path string, parent Group, router *mux.Router) *group {
	return &group{
		parent:      parent,
		groupPath:   path,
		router:      router,
		children:    make([]Group, 0),
		before:      make([]Handler, 0),
		after:       make([]Handler, 0),
		afterAlways: make([]func(*Context), 0),
	}
}

// Use adds the handlers that will run before each request. Note that if a
// handler returns an error, the next handlers won't run.
func (g *group) UseBefore(handlers ...Handler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.before = append(g.before, handlers...)
}

// UseAfter adds the handlers that will run after each request if it is not cancelled.
func (g *group) UseAfter(handlers ...Handler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.after = append(g.after, handlers...)
}

// UseAfterAlways adds the handlers that will run after every request.
func (g *group) UseAfterAlways(handlers ...func(*Context)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.afterAlways = append(g.afterAlways, handlers...)
}

// Group creates a new group to handle http requests.
func (g *group) Group(relativePath string) Group {
	g.mu.Lock()
	defer g.mu.Unlock()

	ng := newGroup(relativePath, g, g.router)
	g.children = append(g.children, ng)
	return ng
}

// Get handles a GET request.
func (g *group) Get(relativePath string, handler Handler) {
	path := joinPaths(g.path(), relativePath)
	h := g.handlerToHttpHandler(handler)
	g.router.HandleFunc(path, h).Methods(http.MethodGet)
}

// Get handles a POST request.
func (g *group) Post(relativePath string, handler Handler) {
	path := joinPaths(g.path(), relativePath)
	h := g.handlerToHttpHandler(handler)
	g.router.HandleFunc(path, h).Methods(http.MethodPost)
}

// Get handles a PUT request.
func (g *group) Put(relativePath string, handler Handler) {
	path := joinPaths(g.path(), relativePath)
	h := g.handlerToHttpHandler(handler)
	g.router.HandleFunc(path, h).Methods(http.MethodPut)
}

// Get handles a DELETE request.
func (g *group) Delete(relativePath string, handler Handler) {
	path := joinPaths(g.path(), relativePath)
	h := g.handlerToHttpHandler(handler)
	g.router.HandleFunc(path, h).Methods(http.MethodDelete)
}

func (g *group) path() string {
	result := ""
	if g.parent != nil {
		result += g.parent.path()
	}

	return joinPaths(result, g.groupPath)
}

func (g *group) getBeforeHandlers() []Handler {
	result := []Handler{}
	if g.parent != nil {
		parentHandlers := g.parent.getBeforeHandlers()
		result = append(result, parentHandlers...)
	}

	result = append(result, g.before...)
	return result
}

func (g *group) getAfterHandlers() []Handler {
	result := []Handler{}
	if g.parent != nil {
		parentHandlers := g.parent.getAfterHandlers()
		result = append(result, parentHandlers...)
	}

	result = append(result, g.after...)
	return result
}

func (g *group) getAfterAlwaysHandlers() []func(*Context) {
	result := []func(*Context){}
	if g.parent != nil {
		parentHandlers := g.parent.getAfterAlwaysHandlers()
		result = append(result, parentHandlers...)
	}

	result = append(result, g.afterAlways...)
	return result
}

func (g *group) handlerToHttpHandler(handler Handler) http.HandlerFunc {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Create request handlers chain.
	before := []Handler{}
	after := []Handler{}
	afterAlways := []func(*Context){}

	if g.parent != nil {
		before = g.parent.getBeforeHandlers()
		after = g.parent.getAfterHandlers()
		afterAlways = g.parent.getAfterAlwaysHandlers()
	}

	before = append(before, g.before...)
	after = append(after, g.after...)
	afterAlways = append(afterAlways, g.afterAlways...)

	// Return the HTTP handlers.
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)

		defer func() {
			// Handle panics.
			if rec := recover(); rec != nil {
				var err error
				switch r := rec.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%s", r)
				}

				ctx.Cancel(err)
			}

			// Run the middlewares that always will run after the request.
			for _, h := range afterAlways {
				h(ctx)
			}

			// Close the response. It writes all data on it.
			err := ctx.closeResponse()

			// Handle any error on closing the response.
			if err != nil {
				// Log the error.
				ctx.Logger().With("error", err.Error()).Error("invalid response content")

				// Set the 500 status code.
				ctx.w.WriteHeader(http.StatusInternalServerError)

				// Set the response body with the internal error.
				internalErr := NewServerError(InternalServerErrorCode, err)
				data, err := m.Marshal(internalErr)
				if err != nil {
					ctx.Logger().With("error", err.Error()).Error("the internal server error is invalid")
				} else {
					ctx.w.Write(data)
				}
			}
		}()

		var reqErr error

		// Run before middlewares.
		for _, h := range before {
			reqErr = h(ctx)
			if reqErr != nil {
				ctx.Cancel(reqErr)
				break
			}
		}

		// Run the request handler.
		if reqErr == nil {
			reqErr = handler(ctx)
		}

		// Run after middlewares.
		if reqErr == nil {
			for _, h := range after {
				reqErr = h(ctx)
				if reqErr != nil {
					ctx.Cancel(reqErr)
					break
				}
			}
		}
	}
}

func joinPaths(path1 string, path2 string) string {
	u, err := url.Parse(path1)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, path2)
	return u.String()
}
