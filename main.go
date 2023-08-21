package main

import (
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	r := gin.Default()
	r.Any("/*any", func(c *gin.Context) {
		// 检查请求头中的"Accept"字段
		acceptHeader := c.GetHeader("Content-Type")
		if strings.Contains(acceptHeader, "application/json") {
			c.Header("Content-Type", "application/json")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
