package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Frame"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if isDevelopmentEnv() {

	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Any("/*url", Frame.TinyRouteHandler)
	r.Run("127.0.0.1:8081") //TODO 配置文件读取

}

func isDevelopmentEnv() bool {
	env := os.Getenv("TPV2_ENV")
	if env == "development" {
		return true
	} else {
		return false
	}
}
