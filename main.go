package main

import (
	"Gin/GinLike"
	"net/http"
)

func main() {
	r := GinLike.Default()
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
