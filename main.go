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
		"/reproxy-openai/*any":    "https://api.openai.com",
		"/reproxy-workergpt/*any": "https://workergpt.cn",
		"/reproxy-eqing/*any":     "https://next.eqing.tech",
		// "/reproxy-miyadns/*any":   "https://www.miyadns.com",
		"/reproxy-chatanywhere/*any": "https://api.chatanywhere.cn",
		"/reproxy-lbbai/*any":        "https://postapi.lbbai.cc",
		"/reproxy-mixerbox/*any":     "https://chatai.mixerbox.com",
		"/reproxy-nnwlink/*any":      "https://api.mmw1984.link",
		"/reproxy-openkey/*any":      "https://openkey.cloud",
	}

	for routePath, targetURL := range proxySites {
		target, _ := url.Parse(targetURL)
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
			},
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			ModifyResponse: func(res *http.Response) error {
				// 添加CORS响应头
				res.Header.Set("Access-Control-Allow-Origin", "*")
				res.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				res.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
				// res.Header.Set("Content-Type", "text/event-stream")
				res.Header.Set("Cache-Control", "no-cache")
				res.Header.Set("Connection", "keep-alive")
				return nil
			},
		}

		r.Any(routePath, func(c *gin.Context) {
			// 附加原始请求路径到目标 URL
			c.Request.URL.Path = target.Path + c.Param("any")

			// Use a new http client with check for redirect
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Transport: proxy.Transport,
			}

			proxy.Transport = client.Transport
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
