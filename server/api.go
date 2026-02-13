package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lfcypo/scrkun/model"
	"github.com/lfcypo/scrkun/store"
)

func SetupAPI(router *gin.Engine) {
	router.Match([]string{"GET", "POST"}, "/api/data", GetAllData)
}

func GetAllData(c *gin.Context) {
	db := store.GetDatabase()

	data := new([]model.Event)
	db.Find(&data)

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
