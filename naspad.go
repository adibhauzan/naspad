package naspad

import "net/http"

// Driver handles routing and implements the http.Handler interface
type Driver struct {
	RGroup
	routes map[string]map[string]HandlerFunc
}

// NewDriver creates a new Driver instance
func NewDriver() *Driver {
	return &Driver{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

// AddRoute registers a route for a specific method and path
func (d *Driver) AddRoute(method, path string, handler HandlerFunc) {
	if _, exists := d.routes[path]; !exists {
		d.routes[path] = make(map[string]HandlerFunc)
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
	handler(NewContext(w, r))
}
