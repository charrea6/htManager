package web

import (
	"htManager/internal/devices"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitWebServer(devices devices.Devices) {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	initAPI(r.Group("/api"), devices)

	initFrontend(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
