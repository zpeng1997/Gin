package main

import (
	"Gin/GinLike"
	"net/http"
)

func main() {
	r := GinLike.Default()
	r.GET("/", func(conn *GinLike.Context){
		conn.HTML(http.StatusOK, "<h1>Hello Gopher!</h1>")
	})
	r.GET("/book", func(conn *GinLike.Context){
		conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Query("name"), conn.Path)
	})
	r.GET("/book/:name", func(conn *GinLike.Context){
		conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Param("name"), conn.Path)
	})
	r.GET("/asset/*filepath", func(conn *GinLike.Context){
		conn.JSON(http.StatusOK, GinLike.H{
			"filepath": conn.Param("filepath"),
		})
	})
	r.POST("/login", func(conn *GinLike.Context) {
		conn.JSON(http.StatusOK, GinLike.H{
			"username": conn.PostForm("username"),
			"password": conn.PostForm("password"),
		})
	})
	r.Run(":9000")
}
