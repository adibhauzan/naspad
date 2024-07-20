package naspad

import (
	"net/http"
	"regexp"
)

// HandlerFunc defines the function signature for handlers
type HandlerFunc func(*Context)

// HandlerPipeline is a slice of HandlerFunc
type HandlerPipeline []HandlerFunc

// Router interface
type Router interface {
	RoutesInterface
}

// RoutesInterface provides methods to handle routes
type RoutesInterface interface {
	Use(middleware ...HandlerFunc) RoutesInterface
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

func (r *RGroup) Use(middleware ...HandlerFunc) RoutesInterface {
	r.Handlers = append(r.Handlers, middleware...)
	return r.returnObj()
}

// BaseRoutePath returns the base route path for the group
func (group *RGroup) BaseRoutePath() string {
	return group.baseRoutePath
}

func (r *RGroup) handle(method, relativePath string, handlers ...HandlerFunc) RoutesInterface {
	// Validate HTTP method using regex
	var regEnLetter = regexp.MustCompile("^[A-Z]+$")
	if !regEnLetter.MatchString(method) {
		panic("HTTP method " + method + " is not valid")
	}

	// Calculate the absolute path for the route
	absolutePath := r.calculateAbsolutePath(relativePath)

	// Combine all handlers into a single HandlerFunc
	combinedHandlers := r.combineHandlers(handlers...)

	// Register the route with the driver
	r.driver.AddRoute(method, absolutePath, combinedHandlers)

	return r
}

func (r *RGroup) Handle(method, relativePath string, handlers ...HandlerFunc) RoutesInterface {
	// Validate the HTTP method
	var regEnLetter = regexp.MustCompile("^[A-Z]+$")
	if !regEnLetter.MatchString(method) {
		panic("HTTP method " + method + " is not valid")
	}

	// Delegate to the handle method to register the route
	return r.handle(method, relativePath, handlers...)
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

// combineHandlers combines multiple handlers into a single HandlerFunc
func (r *RGroup) combineHandlers(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
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
