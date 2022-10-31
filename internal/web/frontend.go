package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"os"
)

//go:embed frontend/build/*
var content embed.FS

func initFrontend(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		rootFS, _ := fs.Sub(content, "frontend/build")

		if _, err := rootFS.Open(c.Request.URL.Path); os.IsNotExist(err) {
			c.FileFromFS("index.html", http.FS(rootFS))
		}
		c.FileFromFS(c.Request.URL.Path, http.FS(rootFS))
	})
}
