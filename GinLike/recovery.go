package GinLike

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)


// 中间件Recovery()
func Recovery() HandlerFunc{
	return func(c *Context){
		defer func(){
			// 捕获panic, 保证不宕机
			if err := recover(); err != nil{
				message := fmt.Sprintf("%s", err)
				// trace 用来获取panic信息
				log.Printf("%s\n\n", trace(message))
				log.Fatal(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

// 这块貌似有点难懂
// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}