package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	proxyURL := os.Getenv("Proxy_URL")
	target, _ := url.Parse(proxyURL)

	r := gin.Default()

	addTransparentProxyRoute(r, "/chatanywhere/*path", target)
	addTransparentProxyRoute(r, "/ohmygpt/*path", target)
	r.NoRoute(func(c *gin.Context) {
		proxy(c.Writer, c.Request, target)
	})

	r.Run()
}

// 添加透明代理路由
func addTransparentProxyRoute(r *gin.Engine, routePath string, target *url.URL) {
	r.Any(routePath, func(c *gin.Context) {
		proxy(c.Writer, c.Request, target)
	})
}

// 透明代理
func proxy(w http.ResponseWriter, r *http.Request, target *url.URL) {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
