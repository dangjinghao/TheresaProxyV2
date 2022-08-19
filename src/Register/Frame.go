package Register

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

var MiddlewareCore []*gin.HandlerFunc

func AddMiddlewareFunc(middlewareFunc gin.HandlerFunc) {

	MiddlewareCore = append(MiddlewareCore, &middlewareFunc)
	return
}

func GetPluginLogger(name string) *logrus.Entry {
	pluginLogger := logrus.WithFields(logrus.Fields{
		"name": name,
	})
	return pluginLogger
}

// 传入不带后缀的插件文件名
func GetPluginConfig(name string) (*os.File, error) {
	return os.Open("config/" + name + ".json")

}
func init() {
	MiddlewareCore = make([]*gin.HandlerFunc, 0)
}
