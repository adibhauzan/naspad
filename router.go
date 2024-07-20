package naspad

import (
	"net/http"
	"regexp"
)

// Router interface
type Router interface {
	RoutesInterface
}

// RoutesInterface provides methods to handle routes
type RoutesInterface interface {
	Handle(method, path string, handler ...HandlerFunc) RoutesInterface
	GET(path string, handler HandlerFunc) RoutesInterface
	POST(path string, handler HandlerFunc) RoutesInterface
}

// RGroup represents a group of routes
type RGroup struct {
	Handlers       HandlerPipeline
	root           bool
	baseRoutePath  string
	driver         *Driver
	routes         map[string]map[string]HandlerFunc
}

// Ensure RGroup implements Router interface
var _ Router = (*RGroup)(nil)

// NewRGroup creates a new RGroup with a given base route path
func NewRGroup(baseRoutePath string, driver *Driver) *RGroup {
	return &RGroup{
		baseRoutePath: baseRoutePath,
		routes:        make(map[string]map[string]HandlerFunc),
		driver:        driver,
	}
}

// BaseRoutePath returns the base route path for the group
func (group *RGroup) BaseRoutePath() string {
	return group.baseRoutePath
}

// Handle registers a handler for a specific HTTP method and path
func (r *RGroup) Handle(method, relativePath string, handlers ...HandlerFunc) RoutesInterface {
	var regEnLetter = regexp.MustCompile("^[A-Z]+$")
	if !regEnLetter.MatchString(method) {
		panic("http method " + method + " is not valid")
	}
	absolutePath := r.calculateAbsolutePath(relativePath)
	combinedHandlers := r.combineHandlers(handlers...)
	r.driver.AddRoute(method, absolutePath, combinedHandlers)
	return r
}

// GET registers a handler for GET requests
func (r *RGroup) GET(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodGet, path, handler)
}

// POST registers a handler for POST requests
func (r *RGroup) POST(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodPost, path, handler)
}

// DELETE registers a handler for DELETE requests
func (r *RGroup) DELETE(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodDelete, path, handler)
}

// PATCH registers a handler for PATCH requests
func (r *RGroup) PATCH(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodPatch, path, handler)
}

// PUT registers a handler for PUT requests
func (r *RGroup) PUT(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodPut, path, handler)
}

// OPTIONS registers a handler for OPTIONS requests
func (r *RGroup) OPTIONS(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodOptions, path, handler)
}

// HEAD registers a handler for HEAD requests
func (r *RGroup) HEAD(path string, handler HandlerFunc) RoutesInterface {
	return r.Handle(http.MethodHead, path, handler)
}

// calculateAbsolutePath calculates the absolute path from the base route path and relative path
func (r *RGroup) calculateAbsolutePath(relativePath string) string {
	if r.baseRoutePath == "" {
		return relativePath
	}
	return r.baseRoutePath + relativePath
}

// combineHandlers combines multiple handlers into a single handler
func (r *RGroup) combineHandlers(handlers ...HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req)
		for _, handler := range handlers {
			handler(ctx)
		}
	}
}

// returnObj returns the appropriate RoutesInterface implementation
func (group *RGroup) returnObj() RoutesInterface {
	if group.root {
		return group.driver
	}
	return group
}
