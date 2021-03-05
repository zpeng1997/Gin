package main

import (
	"Gin/GinLike"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type student struct{
	Name string
	Age int8
}

func FormatAsDate(t time.Time) string{
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%2d-%2d",year, month, day)
}


func onlyFunc() GinLike.HandlerFunc{
	return func(c *GinLike.Context){
		t := time.Now()
		//log.Fatal(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := GinLike.Default()

	// 静态注册文件
	//r.Static("/assets", "/usr/geektutu/blog/static")
	//用户访问localhost:9999/assets/js/geektutu.js，最终返回/usr/geektutu/blog/static/js/geektutu.js

	// 全局的 中间件, Global Middleware
	r.Use(GinLike.Logger())
	r.Use(GinLike.Recovery())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	//加载到内存
	r.LoadHTMLGlob("templates/*")
	// 这块加载文件的还是有点没理解.
	r.Static("/assets", "./static")
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 20}
	r.GET("/", func(c *GinLike.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *GinLike.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", GinLike.H{
			"title":"GinLike",
			"StuArr":[2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *GinLike.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", GinLike.H{
			"title": "GinLike",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	// test Recovery()
	// search for log and find out bugs
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *GinLike.Context) {
		names := []string{"gin like music"}
		c.String(http.StatusOK, names[100])
	})


	//v1 := r.Group("/v1")
	//{
	//	v1.GET("/book", func(conn *GinLike.Context) {
	//		conn.HTML(http.StatusOK, "<h1>Hello Gopher!</h1>")
	//	})
	//	v1.GET("/hello", func(conn *GinLike.Context) {
	//		conn.String(http.StatusOK, "hello %s, you're at %s\n", conn.Param("name"), conn.Path)
	//	})
	//}
	//v2 := r.Group("/v2")
	//// 组的中间件, Group Middleware
	//v2.Use(onlyFunc())
	//{
	//	v2.GET("/asset/*filepath", func(conn *GinLike.Context) {
	//		conn.JSON(http.StatusOK, GinLike.H{
	//			"filepath": conn.Param("filepath"),
	//		})
	//	})
	//	v2.POST("/login", func(conn *GinLike.Context) {
	//		conn.JSON(http.StatusOK, GinLike.H{
	//			"username": conn.PostForm("username"),
	//			"password": conn.PostForm("password"),
	//		})
	//	})
	//}
	r.Run(":9000")
}
