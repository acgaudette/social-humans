package main

import (
	"./app"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	routes *node
}

// Create new router with index handler
func NewRouter(index app.Handler) *Router {
	handlers := make(methodHandlers)
	handlers[http.MethodGet] = index

	// Initialize trie
	routes := &node{
		split:    "",
		children: []*node{},
		handlers: handlers,
	}

	return &Router{routes}
}

// Satisfies multiplexer interface
func (this *Router) ServeHTTP(out http.ResponseWriter, in *http.Request) {
	// Get path tokens
	tokens := strings.Split(in.URL.Path, "/")

	// Special case for root path
	if tokens[1] == "" {
		tokens = tokens[1:]
	}

	// Execute handler for route
	err := this.routes.eval(tokens, out, in)

	if err != nil {
		log.Printf("%s", err)
	}
}

// Store a route (path to method handler) in the router
func (this *Router) Handle(method, route string, handler app.Handler) error {
	// Invalid route
	if route[0] != '/' {
		return fmt.Errorf("invalid route \"%s\"", route)
	}

	// Get path tokens
	tokens := strings.Split(route, "/")

	// Special case for root path
	if tokens[1] == "" {
		tokens = tokens[1:]
	}

	// Add new route
	this.routes.append(method, tokens, handler)
	return nil
}

// GET helper
func (this *Router) GET(route string, handler app.Handler) error {
	return this.Handle(http.MethodGet, route, handler)
}

// POST helper
func (this *Router) POST(route string, handler app.Handler) error {
	return this.Handle(http.MethodPost, route, handler)
}

/* Route trie implementation */

type methodHandlers map[string]app.Handler

type node struct {
	split    string
	children []*node
	handlers methodHandlers
}

// Append a handler to the trie
func (this *node) append(method string, tokens []string, handler app.Handler) {
	// Pop the first token
	tokens = tokens[1:]

	// Base case
	if len(tokens) == 0 {
		// Add first handler
		if this.handlers == nil {
			this.handlers = make(methodHandlers)
		}

		// Otherwise, add handler for new method
		this.handlers[method] = handler
		return
	}

	token := tokens[0]

	// If a child exists with the token, call recursively
	if child := this.getMatch(token); child != nil {
		child.append(method, tokens, handler)
		return
	}

	// Otherwise, create new node
	next := &node{
		split:    token,
		children: []*node{},
		handlers: nil,
	}

	// Add child and call recursively
	this.add(next)
	next.append(method, tokens, handler)
}

// Recursively look for and execute a handler
func (this *node) eval(
	tokens []string, out http.ResponseWriter, in *http.Request,
) error {
	tokens = tokens[1:]

	// Base case
	if len(tokens) == 0 {
		// If handler exists for method, execute
		if handler, ok := this.handlers[in.Method]; ok {
			app.Handle(handler, out, in)
			return nil
		}

		// Error 403
		http.Error(
			out,
			http.StatusText(http.StatusForbidden),
			http.StatusForbidden,
		)

		return fmt.Errorf("method not found for route \"%s\"", in.URL.Path)
	}

	token := tokens[0]

	// If a child exists with the token, call recursively
	if child := this.get(token); child != nil {
		return child.eval(tokens, out, in)
	}

	// Otherwise, error
	http.NotFound(out, in)
	return fmt.Errorf("route not found for \"%s\"", in.URL.Path)
}

// Add child
func (this *node) add(next *node) {
	this.children = append(this.children, next)
}

// Get child with matching split or wildcard
func (this *node) get(split string) *node {
	var wildcard, target *node = nil, nil

	for _, child := range this.children {
		// Check for wildcard match
		if child.split == "*" {
			wildcard = child

			// Break if both have been found
			if target != nil {
				break
			}
		}

		// Check for direct match
		if child.split == split {
			target = child

			// Break if both have been found
			if wildcard != nil {
				break
			}
		}
	}

	// Return direct match over wildcard
	if target != nil {
		return target
	}

	return wildcard
}

// Get child explicitly with a matching split
func (this *node) getMatch(split string) *node {
	for _, child := range this.children {
		if child.split == split {
			return child
		}
	}

	return nil
}
