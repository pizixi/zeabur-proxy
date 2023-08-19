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
	proxy := httputil.NewSingleHostReverseProxy(target)

	r := gin.Default()

	addTransparentReverseProxyRoute(r, "/chatanywhere/*path", "https://api.chatanywhere.cn")
	addTransparentReverseProxyRoute(r, "/ohmygpt/*path", "https://api.ohmygpt.com")
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run()
}

// 添加透明反向代理路由
func addTransparentReverseProxyRoute(r *gin.Engine, routePath string, targetURL string) {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Server")
		return nil
	}

	r.Any(routePath, func(c *gin.Context) {
		c.Request.URL.Path = c.Param("path")
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
