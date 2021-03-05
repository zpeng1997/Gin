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
	"html/template"
	"net/http"
	"path"
	"strings"
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
		// html render
		htmlTemplates *template.Template
		funcMap template.FuncMap
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

// add middlewares handle functions...
func (routerGroup *RouterGroup) Use(middlewares ...HandlerFunc){
	routerGroup.middlewares = append(routerGroup.middlewares, middlewares...)
}

func (routerGroup *RouterGroup) Run(address string)(err error){
	return http.ListenAndServe(address, routerGroup.engine)
}

// 一种是通过addRoute动态注册的函数.
// 另一种是中间件, 在ServeHTTP事先写好, 对所有满足的路由都全局执行. 通过上下文Context传过去
// 最后在router.handler()一起执行.
func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request){
	var middlewares []HandlerFunc
	// 按照组的方式全局覆盖.
	for _, group := range engine.groups{
		// 输入的URL是否 触碰到 prefix
		if strings.HasPrefix(r.URL.Path, group.prefix){}
		middlewares = append(middlewares, group.middlewares...)
	}
	conn := newContext(w, r)
	// 区分middlewares的函数 和 engine 继承 ServeHTTP 之后的 GET, PUT
	conn.handlers = middlewares
	conn.engine = engine
	engine.router.handle(conn)
}

// 得到文件路径
// 并把文件发送过去
func (routerGroup *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc{
	absolutePath := path.Join(routerGroup.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context){
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil{
			c.Status(http.StatusNotFound)
			return
		}
		// 发送文件过去
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 把对应的路径和文件注册上去.
// 当然这是静态文件
func (routerGroup *RouterGroup) Static(relativePath string, root string){
	handler := routerGroup.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	routerGroup.GET(urlPattern, handler)
}

// html render
func (engine *Engine) SetFuncMap(funcMap template.FuncMap){
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string){
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}