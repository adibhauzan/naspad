package naspad

import (
	"encoding/json"
	"net/http"
)

// Context holds the request and response objects
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

// NewContext creates a new Context instance
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

// JSON writes a JSON response
func (c *Context) JSON(status int, v interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	json.NewEncoder(c.Writer).Encode(v)
}

// QueryParam retrieves a query parameter from the URL
func (c *Context) QueryParam(name string) string {
	return c.Request.URL.Query().Get(name)
}

// FormValue retrieves a form value from the request
func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

// SetHeader sets a header for the response
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Error writes an error response
func (c *Context) Error(status int, message string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	json.NewEncoder(c.Writer).Encode(map[string]string{"error": message})
}
