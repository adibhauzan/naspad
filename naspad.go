package naspad

import "net/http"

type HandlerFunc func(*Context)

type HandlerPipeline []HandlerFunc

type Driver struct {
	RGroup
	routes map[string]map[string]http.HandlerFunc
}

// NewDriver creates a new Driver instance
func NewDriver() *Driver {
	return &Driver{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// AddRoute registers a route for a specific method and path
func (d *Driver) AddRoute(method, path string, handler http.HandlerFunc) {
	if _, exists := d.routes[path]; !exists {
		d.routes[path] = make(map[string]http.HandlerFunc)
	}
	d.routes[path][method] = handler
}

// ServeHTTP implements the http.Handler interface for Driver
func (d *Driver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodRoutes, exists := d.routes[r.URL.Path]
	if !exists {
		http.NotFound(w, r)
		return
	}
	handler, exists := methodRoutes[r.Method]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	handler(w, r)
}