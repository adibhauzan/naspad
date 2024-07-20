package naspad

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request   *http.Request
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

func (c *Context) QueryParam(name string) string {
    return c.Request.URL.Query().Get(name)
}

func (c *Context) FormValue(name string) string {
    return c.Request.FormValue(name)
}

func (c *Context) SetHeader(key, value string) {
    c.Writer.Header().Set(key, value)
}

func (c *Context) Error(status int, message string) {
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.WriteHeader(status)
    json.NewEncoder(c.Writer).Encode(map[string]string{"error": message})
}



