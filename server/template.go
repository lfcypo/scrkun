package server

import (
	templatelib "html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lfcypo/scrkun/template"
)

func GetTemplate() *templatelib.Template {
	return templatelib.Must(templatelib.New("").ParseFS(template.TemplateFS, "*.html"))
}

func SetupTemplate(router *gin.Engine) {
	templ := GetTemplate()

	router.SetHTMLTemplate(templ)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
