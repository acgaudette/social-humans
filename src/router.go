package main

import (
	"./front"
	"errors"
	"log"
	"net/http"
	"strings"
)

type router struct {
	routes *node
}

type Handler func(http.ResponseWriter, *http.Request)

func newRouter(index Handler) *router {
	handlers := make(methodHandlers)
	handlers[http.MethodGet] = index

	routes := &node{
		split:    "",
		children: []*node{},
		handlers: handlers,
	}

	return &router{routes: routes}
}

func (this *router) ServeHTTP(
	writer http.ResponseWriter, request *http.Request,
) {

	tokens := strings.Split(request.URL.Path, "/")

	if tokens[1] == "" {
		tokens = tokens[1:]
	}

	err := this.routes.eval(tokens, writer, request)

	if err != nil {
		log.Printf("%s", err)
	}
}

func (this *router) handle(
	method, route string, handler Handler,
) error {
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

func (this *router) GET(route string, handler Handler) error {
	return this.handle(http.MethodGet, route, handler)
}

func (this *router) POST(route string, handler Handler) error {
	return this.handle(http.MethodPost, route, handler)
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
	tokens []string, writer http.ResponseWriter, request *http.Request,
) error {
	tokens = tokens[1:]

	if len(tokens) == 0 {
		if handler, ok := this.handlers[request.Method]; ok {
			handler(writer, request)
			return nil
		}

		front.Error403(writer)
		return errors.New("method not found for route")
	}

	token := tokens[0]

	if child := this.get(token); child != nil {
		return child.eval(tokens, writer, request)
	}

	http.NotFound(writer, request)
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
