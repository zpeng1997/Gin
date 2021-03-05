package main

import (
	"Gin/GinLike"
	"net/http"
)

func main() {
	r := GinLike.Default()
	r.GET("/hello", func(conn *GinLike.Context){
		conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Query("name"), conn.Path)
	})
	r.GET("/book", func(conn *GinLike.Context){
		conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Query("name"), conn.Path)
	})
	r.POST("/login", func(conn *GinLike.Context) {
		conn.JSON(http.StatusOK, GinLike.H{
			"username": conn.PostForm("username"),
			"password": conn.PostForm("password"),
		})
	})
	r.Run(":9000")
}
