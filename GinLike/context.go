package GinLike

import (
	"encoding/json"
	"fmt"
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
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w,
		r,
		r.URL.Path,
		r.Method,
		 make(map[string]string, 0),
		0,
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

func (c *Context) HTML(code int, html string){
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		fmt.Println( "HTML Write error: %v", err)
	}
}