package GinLike

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}

type Context struct{
	// origin object
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	// Response info
	Params map[string]string
	StatusCode int
	// middleware functions
	handlers []HandlerFunc
	index int
	//
	engine *Engine
}

// 初始化方式
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

// 这部分逻辑很重要,
// 保证了中间件, 按照顺序"并发"执行
func (c *Context) Next(){
	c.index ++
	s := len(c.handlers)
	for ; c.index < s; c.index ++{
		c.handlers[c.index](c)
	}
}

// Context get params after router.parsePattern
func (c *Context)Param(key string) string{
	value, _ := c.Params[key]
	return value
}

// reference 《Go Web》, then know how to use http api
func (c *Context) PostForm(key string) string{
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string{
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int){
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string){
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, value ...interface{}){
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
	if err != nil {
		fmt.Println( "String Write error: %v", err)
	}
}

func (c *Context) JSON(code int, obj interface{}){
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data[]byte){
	c.Status(code)
	_, err := c.Writer.Write(data)
	if err != nil {
		fmt.Println( "Data Write error: %v", err)
	}
}

func (c *Context) HTML(code int, name string, data interface{}){
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	//_, err := c.Writer.Write([]byte(html))
	//if err != nil {
	//	fmt.Println( "HTML Write error: %v", err)
	//}
	if err:=c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil{
		log.Fatal(500, err.Error())
	}
}