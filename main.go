package main

import(
	"Gin/GinLike"
	"fmt"
	"net/http"
)

func main() {
	r := GinLike.Default()
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	r.GET("/book", func(w http.ResponseWriter, r *http.Request){
		for k, v := range r.Header{
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	r.Run(":9000")
}
