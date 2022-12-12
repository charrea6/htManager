package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed frontend/build/*
var content embed.FS

func initFrontend(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		rootFS, _ := fs.Sub(content, "frontend/build")
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		_, err := rootFS.Open(path)
		if os.IsNotExist(err) {
			c.FileFromFS("/", http.FS(rootFS))
		} else {
			c.FileFromFS(c.Request.URL.Path, http.FS(rootFS))
		}
	})
}
