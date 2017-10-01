package main

import (
	"./front"
	"errors"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	routes *node
}

func NewRouter(index Handler) *Router {
	handlers := make(methodHandlers)
	handlers[http.MethodGet] = index

	routes := &node{
		split:    "",
		children: []*node{},
		handlers: handlers,
	}

	return &Router{routes: routes}
}

func (this *Router) ServeHTTP(out http.ResponseWriter, in *http.Request) {
	tokens := strings.Split(in.URL.Path, "/")

	if tokens[1] == "" {
		tokens = tokens[1:]
	}

	err := this.routes.eval(tokens, out, in)

	if err != nil {
		log.Printf("%s", err)
	}
}

// Store a route in the router
func (this *Router) Handle(method, route string, handler Handler) error {
	if route[0] != '/' {
		return errors.New("invalid route")
	}

	tokens := strings.Split(route, "/")

	if tokens[1] == "" {
		tokens = tokens[1:]
	}

	this.routes.append(method, tokens, handler)
	return nil
}

// GET helper
func (this *Router) GET(route string, handler Handler) error {
	return this.Handle(http.MethodGet, route, handler)
}

// POST helper
func (this *Router) POST(route string, handler Handler) error {
	return this.Handle(http.MethodPost, route, handler)
}

type methodHandlers map[string]Handler

type node struct {
	split    string
	children []*node
	handlers methodHandlers
}

func (this *node) append(method string, tokens []string, handler Handler) {
	tokens = tokens[1:]

	if len(tokens) == 0 {
		if this.handlers == nil {
			this.handlers = make(methodHandlers)
		}

		this.handlers[method] = handler
		return
	}

	token := tokens[0]

	if child := this.get(token); child != nil {
		child.append(method, tokens, handler)
		return
	}

	next := &node{
		split:    token,
		children: []*node{},
		handlers: nil,
	}

	this.add(next)
	next.append(method, tokens, handler)
}

func (this *node) eval(
	tokens []string, out http.ResponseWriter, in *http.Request,
) error {
	tokens = tokens[1:]

	if len(tokens) == 0 {
		if handler, ok := this.handlers[in.Method]; ok {
			handle(handler, out, in)
			return nil
		}

		front.Error403(out)
		return errors.New("method not found for route")
	}

	token := tokens[0]

	if child := this.get(token); child != nil {
		return child.eval(tokens, out, in)
	}

	http.NotFound(out, in)
	return errors.New("route not found")
}

func (this *node) add(next *node) {
	this.children = append(this.children, next)
}

func (this *node) get(split string) *node {
	for _, child := range this.children {
		if child.split == split || child.split == "*" {
			return child
		}
	}

	return nil
}
