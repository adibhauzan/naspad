package naspad

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
)

type Context struct {
	Writer  http.ResponseWriter
	Request   *http.Request

	index    int8
	fullPath string


	// This mutex protects Keys map.
	mu sync.RWMutex

	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any

	// Errors is a list of errors attached to all the handlers/middlewares who used this context.

	// Accepted defines a list of manually accepted formats for content negotiation.
	Accepted []string

	// queryCache caches the query result from c.Request.URL.Query().
	queryCache url.Values

	// formCache caches c.Request.PostForm, which contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	formCache url.Values

	// SameSite allows a server to define a cookie attribute making it impossible for
	// the browser to send this cookie along with cross-site requests.
	sameSite http.SameSite
}


func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

func (c *Context) JSON(status int, v interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	json.NewEncoder(c.Writer).Encode(v)
} 



// func (c *Context) Next() {
// 	c.index++
// 	for c.index < int8(len(c.handlers)) {
// 		c.h[c.index](c)
// 		c.index++
// 	}
// }
