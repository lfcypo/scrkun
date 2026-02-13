package server

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	superRouter := gin.Default()

	SetupTemplate(superRouter)
	SetupAPI(superRouter)

	return superRouter
}
