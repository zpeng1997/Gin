package main

import (
	"Gin/GinLike"
	"log"
	"net/http"
	"time"
)

func onlyFunc() GinLike.HandlerFunc{
	return func(c *GinLike.Context){
		t := time.Now()
		//log.Fatal(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := GinLike.Default()
	// 全局的 中间件, Global Middleware
	r.Use(GinLike.Logger())
	r.GET("/index", func(c *GinLike.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/book", func(conn *GinLike.Context) {
			conn.HTML(http.StatusOK, "<h1>Hello Gopher!</h1>")
		})
		v1.GET("/hello", func(conn *GinLike.Context) {
			conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Param("name"), conn.Path)
		})
	}
	v2 := r.Group("/v2")
	// 组的中间件, Group Middleware
	v2.Use(onlyFunc())
	{
		v2.GET("/asset/*filepath", func(conn *GinLike.Context) {
			conn.JSON(http.StatusOK, GinLike.H{
				"filepath": conn.Param("filepath"),
			})
		})
		v2.POST("/login", func(conn *GinLike.Context) {
			conn.JSON(http.StatusOK, GinLike.H{
				"username": conn.PostForm("username"),
				"password": conn.PostForm("password"),
			})
		})
	}
	r.Run(":9000")
}
