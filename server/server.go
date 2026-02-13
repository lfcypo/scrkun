package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lfcypo/scrkun/logger"
)

var serverLog = logger.New("HTTP Server")

func init() {
	gin.SetMode(gin.ReleaseMode)

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		devNull = os.Stdout
	}

	gin.DefaultWriter = devNull
	gin.DefaultErrorWriter = devNull
}

func StartWebServer(port int) {
	router := SetupRouter()

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		serverLog.Fatalf("启动 HTTP 服务器失败: %v", err)
	}
}
