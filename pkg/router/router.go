package router

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	realtimeforum "real-time-forum"
	"strings"
)

type route struct {
	method  string
	pattern string
	parts   []string
	handler http.Handler
}

type contextKey string

const paramsKey contextKey = "params"

type AppRouter struct {
	routes []route

	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func NewRouter() *AppRouter {
	return &AppRouter{}
}

func (r *AppRouter) Handle(method, pattern string, h http.Handler) {

	for _, rt := range r.routes {
		if rt.method == method && rt.pattern == pattern {
			panic("duplicate route: " + method + " " + pattern)
		}
	}

	r.routes = append(r.routes, route{
		method:  method,
		pattern: pattern,
		parts:   splitPath(pattern),
		handler: h,
	})
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}

func (r *AppRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if strings.HasPrefix(req.URL.Path, "/static/") {

		if strings.HasSuffix(req.URL.Path, "/") {
			slog.Warn("static directory listing attempted", "path", req.URL.Path)
			if r.ErrorHandler != nil {
				r.ErrorHandler(w, req, realtimeforum.ErrNotFound)
				return
			}
			http.NotFound(w, req)
			return
		}

		filePath := "./static" + strings.TrimPrefix(req.URL.Path, "/static")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			slog.Warn("static file not found", "path", req.URL.Path)
			r.ErrorHandler(w, req, realtimeforum.ErrNotFound)
			return
		}

		fs := http.FileServer(http.Dir("./static"))
		http.StripPrefix("/static/", fs).ServeHTTP(w, req)
		return
	}

	if strings.HasPrefix(req.URL.Path, "/uploads/") {

		if strings.HasSuffix(req.URL.Path, "/") {
			slog.Warn("uploads directory listing attempted", "path", req.URL.Path)
			if r.ErrorHandler != nil {
				r.ErrorHandler(w, req, realtimeforum.ErrNotFound)
				return
			}
			http.NotFound(w, req)
			return
		}

		filePath := "./uploads" + strings.TrimPrefix(req.URL.Path, "/uploads")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			slog.Warn("upload file not found", "path", req.URL.Path)
			r.ErrorHandler(w, req, realtimeforum.ErrNotFound)
			return
		}

		fs := http.FileServer(http.Dir("uploads"))
		http.StripPrefix("/uploads/", fs).ServeHTTP(w, req)
		return
	}

	reqParts := splitPath(req.URL.Path)
	var allowedMethods []string

	for _, rt := range r.routes {

		ok, params := match(rt.parts, reqParts)
		if !ok {
			continue
		}

		if req.Method == rt.method {
			ctx := context.WithValue(req.Context(), paramsKey, params)
			rt.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		allowedMethods = append(allowedMethods, rt.method)
	}

	if len(allowedMethods) > 0 {
		w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		slog.Warn("method not allowed",
			"method", req.Method,
			"path", req.URL.Path,
			"allowed", strings.Join(allowedMethods, ", "),
		)
		if r.ErrorHandler != nil {
			r.ErrorHandler(w, req, realtimeforum.ErrMethodNotAllowed)
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Warn("route not found", "method", req.Method, "path", req.URL.Path)
	if r.ErrorHandler != nil {
		r.ErrorHandler(w, req, realtimeforum.ErrNotFound)
		return
	}

	http.NotFound(w, req)
}

func match(routeParts, reqParts []string) (bool, map[string]string) {

	params := make(map[string]string)

	if len(routeParts) > 0 && routeParts[len(routeParts)-1] == "" {
		if len(reqParts) < len(routeParts)-1 {
			return false, nil
		}
		for i := 0; i < len(routeParts)-1; i++ {
			if routeParts[i] != reqParts[i] {
				return false, nil
			}
		}
		return true, params
	}

	if len(routeParts) != len(reqParts) {
		return false, nil
	}

	for i := range routeParts {

		if strings.HasPrefix(routeParts[i], "{") &&
			strings.HasSuffix(routeParts[i], "}") {

			key := routeParts[i][1 : len(routeParts[i])-1]
			params[key] = reqParts[i]
			continue
		}

		if routeParts[i] != reqParts[i] {
			return false, nil
		}
	}

	return true, params
}

func Param(r *http.Request, key string) string {
	params, ok := r.Context().Value(paramsKey).(map[string]string)
	if !ok {
		return ""
	}

	val, exists := params[key]
	if !exists {
		return ""
	}

	return val
}

func Params(r *http.Request) map[string]string {
	params, _ := r.Context().Value(paramsKey).(map[string]string)
	return params
}
