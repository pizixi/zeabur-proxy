package main

import (
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	r := gin.Default()
	r.Any("/*any", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
