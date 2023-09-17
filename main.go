package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	proxySites := map[string]string{
		"/reproxy-workergpt/*any": "https://workergpt.cn",
		"/reproxy-eqing/*any":     "https://next.eqing.tech",
		// "/reproxy-miyadns/*any":   "https://www.miyadns.com",
		"/reproxy-chatanywhere/*any": "https://api.chatanywhere.cn",
		"/reproxy-lbbai/*any":        "https://postapi.lbbai.cc",
		"/reproxy-mixerbox/*any":     "https://chatai.mixerbox.com",
	}

	for routePath, targetURL := range proxySites {
		target, _ := url.Parse(targetURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// 修改代理响应
		proxy.ModifyResponse = func(res *http.Response) error {
			// 添加CORS响应头
			res.Header.Set("Access-Control-Allow-Origin", "*")
			res.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			res.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			res.Header.Set("Content-Type", "text/event-stream")
			res.Header.Set("Cache-Control", "no-cache")
			res.Header.Set("Connection", "keep-alive")
			return nil
		}

		r.Any(routePath, func(c *gin.Context) {
			// 附加原始请求路径到目标 URL
			c.Request.URL.Path = target.Path + c.Param("any")
			// 设置请求的Host头
			c.Request.Host = target.Host

			// 添加一個 HTTP 頭部
			c.Request.Header.Add("X-API-Key", "94E160EA-F20D-0123-D7B9-DBE77FB345EF")

			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// 处理未匹配到的路径
	r.NoRoute(func(c *gin.Context) {
		target, _ := url.Parse(os.Getenv("Proxy_URL"))
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":" + os.Getenv("PORT")) // 在0.0.0.0:8899 上监听并在此处服务
}
