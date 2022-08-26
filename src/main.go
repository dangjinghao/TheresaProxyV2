package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Config"
	"TheresaProxyV2/src/Frame"
	"TheresaProxyV2/src/Register"
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
	Config.Version = "0.1.14"

	logger := Config.NewLoggerWithName("main")
	if filePtr, err := os.Open("tpv2Config/tpv2Config.json"); err != nil {
		fmt.Printf("配置文件读取失败:%s,启用默认设置\n", err.Error())
		configStr, _ := Config.TPV2ConfigFs.ReadFile("config.json")
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
	logger.Infof("当前版本为%q", Config.Version)
	r := gin.Default()

	logger.Debug("加载session中间件")
	store := cookie.NewStore([]byte("proxy-secret"))
	r.Use(sessions.Sessions("proxy-session", store))
	logger.Debug("session中间件加载完成")

	for k, _ := range Register.PluginRoute {
		logger.Info("添加拓展路由" + k)
	}

	logger.Debug("加载拓展中间件")
	for _, v := range Register.MiddlewareCore {
		r.Use(*v)
	}
	logger.Debug("拓展中间件加载完成")

	r.Any("/*url", Frame.TinyRouteHandler)
	logger.Infof("TheresaProxyV2启动，地址:%v", tpv2Config.Addr)

	if err := r.Run(tpv2Config.Addr); err != nil {
		logger.Panicf("TheresaProxyV2启动失败:%v", err)
	}

}
