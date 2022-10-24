package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Group interface {
	// Use adds a handler function in the request handlers chain.
	Use(handlers ...http.HandlerFunc)
	// Group creates a new group to handle http requests.
	Group(relativePath string, handlers ...http.HandlerFunc) Group

	// Get handles a GET request.
	Get(relativePath string, handler http.HandlerFunc)
	// Get handles a POST request.
	Post(relativePath string, handler http.HandlerFunc)
	// Get handles a PUT request.
	Put(relativePath string, handler http.HandlerFunc)
	// Get handles a DELETE request.
	Delete(relativePath string, handler http.HandlerFunc)
}

// group is the group to wrap http handlers.
type group struct {
	g *gin.RouterGroup
}

// newGroup creates a new group instance.
func newGroup(g *gin.RouterGroup) *group {
	return &group{
		g: g,
	}
}

// Use adds a handler function in the request handlers chain.
func (g *group) Use(handlers ...http.HandlerFunc) {
	h := handlersToGinHandlers(handlers...)
	g.g.Use(h...)
}

// Group creates a new group to handle http requests.
func (g *group) Group(relativePath string, handlers ...http.HandlerFunc) Group {
	h := handlersToGinHandlers(handlers...)
	ginG := g.g.Group(relativePath, h...)
	return newGroup(ginG)
}

// Get handles a GET request.
func (g *group) Get(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	g.g.GET(relativePath, h)
}

// Get handles a POST request.
func (g *group) Post(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	g.g.GET(relativePath, h)
}

// Get handles a PUT request.
func (g *group) Put(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	g.g.GET(relativePath, h)
}

// Get handles a DELETE request.
func (g *group) Delete(relativePath string, handler http.HandlerFunc) {
	h := gin.WrapF(handler)
	g.g.GET(relativePath, h)
}
