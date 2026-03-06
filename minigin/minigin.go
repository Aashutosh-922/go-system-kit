package minigin

import (
	"encoding/json"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type route struct {
	method   string
	pattern  string
	segments []string
	handlers []HandlerFunc
}

type Engine struct {
	middlewares []HandlerFunc
	routes      []route
}

func New() *Engine {
	return &Engine{
		routes: make([]route, 0),
	}
}

func (e *Engine) Use(mw ...HandlerFunc) {
	e.middlewares = append(e.middlewares, mw...)
}

func (e *Engine) GET(pattern string, handlers ...HandlerFunc) {
	e.addRoute(http.MethodGet, pattern, handlers...)
}

func (e *Engine) POST(pattern string, handlers ...HandlerFunc) {
	e.addRoute(http.MethodPost, pattern, handlers...)
}

func (e *Engine) PUT(pattern string, handlers ...HandlerFunc) {
	e.addRoute(http.MethodPut, pattern, handlers...)
}

func (e *Engine) DELETE(pattern string, handlers ...HandlerFunc) {
	e.addRoute(http.MethodDelete, pattern, handlers...)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	methodMismatch := false

	for _, rt := range e.routes {
		params, ok := match(rt.segments, r.URL.Path)
		if !ok {
			continue
		}

		if rt.method != r.Method {
			methodMismatch = true
			continue
		}

		handlers := make([]HandlerFunc, 0, len(e.middlewares)+len(rt.handlers))
		handlers = append(handlers, e.middlewares...)
		handlers = append(handlers, rt.handlers...)

		ctx := &Context{
			Writer:   w,
			Request:  r,
			Params:   params,
			handlers: handlers,
			index:    -1,
			values:   make(map[string]any),
		}
		ctx.Next()
		return
	}

	if methodMismatch {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func (e *Engine) addRoute(method, pattern string, handlers ...HandlerFunc) {
	if len(handlers) == 0 {
		panic("minigin: route must have at least one handler")
	}

	e.routes = append(e.routes, route{
		method:   method,
		pattern:  pattern,
		segments: splitPath(pattern),
		handlers: handlers,
	})
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string

	handlers []HandlerFunc
	index    int
	values   map[string]any
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

func (c *Context) AbortWithStatus(status int) {
	c.Status(status)
	c.Abort()
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Set(key string, value any) {
	c.values[key] = value
}

func (c *Context) Get(key string) (any, bool) {
	v, ok := c.values[key]
	return v, ok
}

func (c *Context) Status(status int) {
	c.Writer.WriteHeader(status)
}

func (c *Context) String(status int, value string) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(value))
}

func (c *Context) JSON(status int, value any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	_ = json.NewEncoder(c.Writer).Encode(value)
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func match(patternSegments []string, requestPath string) (map[string]string, bool) {
	requestSegments := splitPath(requestPath)
	params := make(map[string]string)

	patternLen := len(patternSegments)
	requestLen := len(requestSegments)

	if patternLen == 0 && requestLen == 0 {
		return params, true
	}

	for i := 0; i < patternLen; i++ {
		if i >= requestLen {
			return nil, false
		}

		segment := patternSegments[i]
		if strings.HasPrefix(segment, ":") {
			params[strings.TrimPrefix(segment, ":")] = requestSegments[i]
			continue
		}

		if strings.HasPrefix(segment, "*") {
			params[strings.TrimPrefix(segment, "*")] = strings.Join(requestSegments[i:], "/")
			return params, true
		}

		if segment != requestSegments[i] {
			return nil, false
		}
	}

	return params, patternLen == requestLen
}

// Usage
//
// app := minigin.New()
// app.Use(func(c *minigin.Context) {
//     log.Println(c.Request.Method, c.Request.URL.Path)
//     c.Next()
// })
// app.GET("/users/:id", func(c *minigin.Context) {
//     c.JSON(200, map[string]string{"id": c.Param("id")})
// })
// _ = app.Run(":8080")
