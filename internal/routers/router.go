package routers

import (
	"fmt"
	"github.com/LiddleChild/http-from-tcp/internal/http"
	"strings"
)

type node struct {
	absolute string
	subpath  string

	handlers    map[http.Method]http.Handler
	dynamicPath *string

	routes map[string]*node
}

type Router struct {
	root *node
}

func NewRouter() *Router {
	return &Router{
		root: newNode("", ""),
	}
}

func newNode(subpath, absolute string) *node {
	return &node{
		absolute:    absolute,
		subpath:     subpath,
		handlers:    map[http.Method]http.Handler{},
		dynamicPath: nil,
		routes:      map[string]*node{},
	}
}

func (r *Router) getSubPaths(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}

func (r *Router) matchLongestPath(path string, dynamic bool) (*node, int) {
	subpaths := r.getSubPaths(path)

	var (
		index   = 0
		current = r.root
	)

	for index < len(subpaths) {
		subpath := subpaths[index]

		next, ok := current.routes[subpath]
		if ok {
			current = next
			index += 1
		} else if dynamic && current.dynamicPath != nil {
			current = current.routes[*current.dynamicPath]
			index += 1
		} else {
			break
		}
	}

	return current, index
}

func (r *Router) RegisterHandler(method http.Method, absolute string, handler http.Handler) error {
	subpaths := r.getSubPaths(absolute)
	parent, index := r.matchLongestPath(absolute, false)

	for _, subpath := range subpaths[index:] {
		dynamicSubpath := strings.HasPrefix(subpath, ":")

		if dynamicSubpath && parent.dynamicPath != nil && *parent.dynamicPath != subpath {
			return fmt.Errorf("%s conflicts with wildcard %s in %s/%s", absolute, subpath, parent.absolute, *parent.dynamicPath)
		}

		if dynamicSubpath {
			parent.dynamicPath = &subpath
		}

		newNode := newNode(subpath, fmt.Sprintf("%s/%s", parent.absolute, subpath))
		parent.routes[subpath] = newNode
		parent = newNode
	}

	if _, ok := parent.handlers[method]; ok {
		return fmt.Errorf("%s %s handler already exists", method, absolute)
	}

	parent.handlers[method] = handler

	return nil
}

func (r *Router) ParseParam(params map[string]string, path string) {
	subpaths := r.getSubPaths(path)
	parent, index := r.matchLongestPath(path, true)

	if len(subpaths) > index {
		return
	}

	subtemplates := r.getSubPaths(parent.absolute)

	for index, subtemplate := range subtemplates {
		if strings.HasPrefix(subtemplate, ":") {
			params[strings.TrimPrefix(subtemplate, ":")] = subpaths[index]
		}
	}
}

func (r *Router) GetHandler(method http.Method, path string) http.Handler {
	subpaths := r.getSubPaths(path)
	parent, index := r.matchLongestPath(path, true)

	if len(subpaths) > index {
		return nil
	}

	handler, ok := parent.handlers[method]
	if !ok {
		return nil
	}

	return handler
}

func (r *Router) ListRoutes() {
	st := []*node{r.root}

	push := func(n *node) {
		st = append(st, n)
	}

	pop := func() *node {
		top := st[len(st)-1]
		st = st[:len(st)-1]
		return top
	}

	for len(st) > 0 {
		top := pop()

		for method := range top.handlers {
			fmt.Printf("%s\t%s\n", method, top.absolute)
		}

		for _, node := range top.routes {
			push(node)
		}
	}
}
