package main

import (
	"TheresaProxyV2/frame"
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/rawConfig"
	"TheresaProxyV2/register"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

type TPV2Config struct {
	Addr string `json:"addr"`
	ENV  string `json:"env"`
}

var tpv2Config TPV2Config

func main() {
	logger := rawConfig.NewLoggerWithName("main")
	if filePtr, err := os.Open("Config/config.json"); err != nil {
		fmt.Printf("配置文件读取失败:%s,启用默认设置\n", err.Error())
		configStr, _ := rawConfig.TPV2ConfigFs.ReadFile("config.json")
		json.Unmarshal(configStr, &tpv2Config)
	} else {
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&tpv2Config)
		filePtr.Close()
	}

	if tpv2Config.ENV == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if tpv2Config.ENV == "production" {
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else {
		panic("未知的环境变量设置")
	}
	logger.Infof("目前环境为：%s", tpv2Config.ENV)
	logger.Infof("当前版本为%v", rawConfig.Version)
	r := gin.Default()

	logger.Debug("加载session中间件")
	store := cookie.NewStore([]byte("proxy-secret"))
	r.Use(sessions.Sessions("proxy-session", store))
	logger.Debug("session中间件加载完成")

	for k, _ := range register.PluginRoute {
		logger.Info("添加拓展路由" + k)
	}

	logger.Debug("加载拓展中间件")
	for _, v := range register.MiddlewareCore {
		r.Use(*v)
	}
	logger.Debug("拓展中间件加载完成")

	r.Any("/*url", frame.TinyRouteHandler)
	logger.Infof("TheresaProxyV2启动，地址:%v", tpv2Config.Addr)

	if err := r.Run(tpv2Config.Addr); err != nil {
		logger.Panicf("TheresaProxyV2启动失败:%v", err)
	}

}
