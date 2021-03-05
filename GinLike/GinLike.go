// version 1.0

//package GinLike
//
//import(
//	"fmt"
//	"net/http"
//)
//
//// HandlerFunc defines the request handler by GinLike
//type HandlerFunc func(http.ResponseWriter, *http.Request)
//
//// Implement the interface of ServerHTTP
//type Engine struct{
//	router map[string]HandlerFunc
//}
//
//func Default() *Engine {
//	return &Engine{make(map[string]HandlerFunc)}
//}
//
//// add handler to correspondingly route(method + pattern)
//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc){
//	key := method + "-" + pattern
//	engine.router[key] = handler
//}
//
//// add GET method to router
//func (engine *Engine) GET(pattern string, handler HandlerFunc){
//	engine.addRoute("GET", pattern, handler)
//}
//
//// add POST method to router
//func (engine *Engine) POST(pattern string, handler HandlerFunc){
//	engine.addRoute("POST", pattern, handler)
//}
//
//// add DELETE method to router
//func (engine *Engine) DELETE(pattern string, handler HandlerFunc){
//	engine.addRoute("DELETE", pattern, handler)
//}
//
//// add PUT method to router
//func (engine *Engine) PUT(pattern string, handler HandlerFunc){
//	engine.addRoute("PUT", pattern, handler)
//}
//
///*
//some other methods can implement by yourself....
// */
//
//func (engine *Engine)Run(address string)(err error){
//	return http.ListenAndServe(address, engine)
//}
//
//func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request){
//	key := r.Method + "-" + r.URL.Path
//	handler, err := engine.router[key]
//	if !err {
//		w.WriteHeader(http.StatusNotFound)
//		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
//	}else{
//		// what does it do?
//		// depends on the implement of correspondingly method
//		handler(w, r)
//	}
//}


// version 2.0
package GinLike

import (
	"net/http"
)

type HandlerFunc func(*Context)

// 这里的处理感觉并不是很好
// 以后记得优化一下!!
type (
	RouterGroup struct{
	prefix string
	middlewares []HandlerFunc // support middleware
	parent *RouterGroup // support nesting
	engine *Engine // all groups share a engine instance
	}

    Engine struct{
    *RouterGroup // 嵌套的写法, 不用写变量名
	router *router
	groups []*RouterGroup
	}
)

// 显式初始化:
func Default() *Engine {
	engine := &Engine{router: newRouter()}
	// 因为没有变量名, 所以直接用 类型名
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (routerGroup *RouterGroup) Group(prefix string) *RouterGroup {
	engine := routerGroup.engine
	newGroup := &RouterGroup{
		prefix: routerGroup.prefix + prefix,
		parent: routerGroup,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}


func (routerGroup *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = routerGroup.prefix + pattern
	routerGroup.engine.router.addRoute(method, pattern, handler)
}

func (routerGroup *RouterGroup) GET(pattern string, handler HandlerFunc){
	routerGroup.addRoute("GET", pattern, handler)
}

// add POST method to router
func (routerGroup *RouterGroup) POST(pattern string, handler HandlerFunc){
	routerGroup.addRoute("POST", pattern, handler)
}

// add DELETE method to router
func (routerGroup *RouterGroup) DELETE(pattern string, handler HandlerFunc){
	routerGroup.addRoute("DELETE", pattern, handler)
}

// add PUT method to router
func (routerGroup *RouterGroup) PUT(pattern string, handler HandlerFunc){
	routerGroup.addRoute("PUT", pattern, handler)
}

/*
some other methods can implement by yourself....
*/

func (routerGroup *RouterGroup)Run(address string)(err error){
	return http.ListenAndServe(address, routerGroup.engine)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request){
	conn := newContext(w, r)
	engine.router.handle(conn)
}