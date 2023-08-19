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

		// Add custom headers
		c.Request.Header.Add("X-Forwarded-For", c.ClientIP())
		c.Request.Header.Add("X-Real-IP", c.ClientIP())
		c.Request.Header.Add("X-Forwarded-Proto", c.Request.URL.Scheme)
		c.Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
		c.Request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.8")
		c.Request.Header.Set("Accept-Encoding", "gzip, deflate, br")
		c.Request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		c.Request.Header.Set("Sec-Ch-Ua", `"Not/0 A Brand";v="99", "Microsoft Edge";v="115", "Chromium";v="115"`)
		c.Request.Header.Set("Sec-Ch-Ua-Platform", "Windows")
		c.Request.Header.Set("Sec-Fetch-Dest", "document")
		c.Request.Header.Set("Sec-Fetch-Mode", "navigate")
		c.Request.Header.Set("Sec-Fetch-Site", "none")
		c.Request.Header.Set("Upgrade-Insecure-Requests", "1")

		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
