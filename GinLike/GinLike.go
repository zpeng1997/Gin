package GinLike

import(
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler by GinLike
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Implement the interface of ServerHTTP
type Engine struct{
	router map[string]HandlerFunc
}

func Default() *Engine {
	return &Engine{make(map[string]HandlerFunc)}
}

// add handler to correspondingly route(method + pattern)
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc){
	key := method + "-" + pattern
	engine.router[key] = handler
}

// add GET method to router
func (engine *Engine) GET(pattern string, handler HandlerFunc){
	engine.addRoute("GET", pattern, handler)
}

// add POST method to router
func (engine *Engine) POST(pattern string, handler HandlerFunc){
	engine.addRoute("POST", pattern, handler)
}

// add DELETE method to router
func (engine *Engine) DELETE(pattern string, handler HandlerFunc){
	engine.addRoute("DELETE", pattern, handler)
}

// add PUT method to router
func (engine *Engine) PUT(pattern string, handler HandlerFunc){
	engine.addRoute("PUT", pattern, handler)
}

/*
some other methods can implement by yourself....
 */

func (engine *Engine)Run(address string)(err error){
	return http.ListenAndServe(address, engine)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request){
	key := r.Method + "-" + r.URL.Path
	handler, err := engine.router[key]
	if !err {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}else{
		// what does it do?
		// depends on the implement of correspondingly method
		handler(w, r)
	}
}