package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lfcypo/scrkun/logger"
)

var serverLog = logger.New("HTTP Server")

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func StartWebServer(port int) {
	router := SetupRouter()

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		serverLog.Fatalf("启动 HTTP 服务器失败: %v", err)
	}
}
