package naspad

import (
	"net/http"
)

var methods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

type Router struct {
	Routes
	handlers    map[string]map[string]HandlerFunc
	middlewares []MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]map[string]HandlerFunc)}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	if handlers, ok := r.handlers[path]; ok {
		if handler, ok := handlers[method]; ok {
			c := NewContext(w, req)
			for _, middleware := range r.middlewares {
				middleware(c)
			}
			handler(c)
			return
		}
	}
	http.NotFound(w, req)
}

type Routes interface {
	Use(...HandlerFunc) Routes

	Handle(string, string, ...HandlerFunc) Routes
	Any(string, ...HandlerFunc) Routes
	GET(string, ...HandlerFunc) Routes
	POST(string, ...HandlerFunc) Routes
	DELETE(string, ...HandlerFunc) Routes
	PATCH(string, ...HandlerFunc) Routes
	PUT(string, ...HandlerFunc) Routes
	OPTIONS(string, ...HandlerFunc) Routes
	HEAD(string, ...HandlerFunc) Routes
	Match([]string, string, ...HandlerFunc) Routes

	StaticFile(string, string) Routes
	StaticFileFS(string, string, http.FileSystem) Routes
	Static(string, string) Routes
	StaticFS(string, http.FileSystem) Routes
}

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

func (r *Router) Handle(method, path string, handlers ...HandlerFunc) Routes {
	if r.handlers[path] == nil {
		r.handlers[path] = make(map[string]HandlerFunc)
	}
	// Chain middlewares and handler functions
	var handler HandlerFunc
	if len(handlers) > 1 {
		handler = handlers[len(handlers)-1]
		for i := len(handlers) - 2; i >= 0; i-- {
			h := handlers[i]
			handler = chainMiddleware(h, handler)
		}
	} else if len(handlers) == 1 {
		handler = handlers[0]
	}
	r.handlers[path][method] = handler
	return r
}

func chainMiddleware(mw, next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		mw(c)
		next(c)
	}
}

func (r *Router) Use(middleware ...HandlerFunc) Routes {
	for _, mw := range middleware {
		r.middlewares = append(r.middlewares, MiddlewareFunc(mw))
	}
	return r
}

func (r *RouterGroup) Use(middleware ...HandlerFunc) Routes {
	r.Handlers = append(r.Handlers, middleware...)
	return r.returnObj()
}

func (r *Router) Any(path string, handlers ...HandlerFunc) Routes {
	for _, method := range methods {
		r.Handle(method, path, handlers...)
	}
	return r
}

func (r *Router) GET(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodGet, path, handlers...)
}

func (r *Router) POST(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodPost, path, handlers...)
}

func (r *Router) DELETE(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodDelete, path, handlers...)
}

func (r *Router) PATCH(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodPatch, path, handlers...)
}

func (r *Router) PUT(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodPut, path, handlers...)
}

func (r *Router) OPTIONS(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodOptions, path, handlers...)
}

func (r *Router) HEAD(path string, handlers ...HandlerFunc) Routes {
	return r.Handle(http.MethodHead, path, handlers...)
}

func (r *Router) Match(methods []string, path string, handlers ...HandlerFunc) Routes {
	for _, method := range methods {
		r.Handle(method, path, handlers...)
	}
	return r
}

func (r *Router) StaticFile(path, filepath string) Routes {
	return r.Handle(http.MethodGet, path, func(c *Context) {
		http.ServeFile(c.Writer, c.Request, filepath)
	})
}

func (r *Router) StaticFileFS(path, filepath string, fs http.FileSystem) Routes {
	return r.Handle(http.MethodGet, path, func(c *Context) {
		http.FileServer(fs).ServeHTTP(c.Writer, c.Request)
	})
}

func (r *Router) Static(path, root string) Routes {
	fs := http.Dir(root)
	fileServer := http.FileServer(fs)
	return r.Handle(http.MethodGet, path, func(c *Context) {
		http.StripPrefix(path, fileServer).ServeHTTP(c.Writer, c.Request)
	})
}

func (r *Router) StaticFS(path string, fs http.FileSystem) Routes {
	fileServer := http.FileServer(fs)
	return r.Handle(http.MethodGet, path, func(c *Context) {
		http.StripPrefix(path, fileServer).ServeHTTP(c.Writer, c.Request)
	})
}

func (r *RouterGroup) returnObj() Routes {
	if r.root {
		return r.engine
	}
	return r
}
