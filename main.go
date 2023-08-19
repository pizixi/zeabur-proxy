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
		SetCORS(c)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

// SetCORS SetCORS
func SetCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, x-requested-with")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, x-requested-with")
	// c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	// c.Writer.Header().Set("Content-Type", "application/json")
}
