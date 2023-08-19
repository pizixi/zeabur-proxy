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

	addReverseProxyRoute(r, "/chatanywhere/*path", "https://api.chatanywhere.cn")
	addReverseProxyRoute(r, "/ohmygpt/*path", "https://api.ohmygpt.com")
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run()
}

// 添加反向代理路由
func addReverseProxyRoute(r *gin.Engine, routePath string, targetURL string) {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

	r.Any(routePath, func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")

		// Add X-Forwarded-For header
		c.Request.Header.Add("X-Forwarded-For", c.ClientIP())
		c.Request.Header.Add("X-Real-IP", c.ClientIP())
		c.Request.Header.Add("X-Forwarded-Proto", c.Request.URL.Scheme)

		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
