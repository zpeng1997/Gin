package GinLike

import (
	"log"
	"time"
)

/*
中间件可以给框架提供无限的扩展能力，应用在分组上，可以使得分组控制的收益更为明显，而不是共享相同的路由前缀这么简单。
例如/admin的分组，可以应用鉴权中间件；/分组应用日志中间件，/是默认的最顶层的分组，
也就意味着给所有的路由，即整个框架增加了记录日志的能力。
 */

/*
中间件是应用在RouterGroup上的，应用在最顶层的 Group，相当于作用于全局，所有的请求都会被中间件处理。
那为什么不作用在每一条路由规则上呢？
作用在某条路由规则，那还不如用户直接在 Handler 中调用直观。
只作用在某条路由规则的功能通用性太差，不适合定义为中间件。
 */

func Logger() HandlerFunc{
	return func(c *Context){
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}