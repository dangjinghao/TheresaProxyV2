package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Config"
	"TheresaProxyV2/src/Frame"
	"TheresaProxyV2/src/Register"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

type TPV2Config struct {
	Addr string `json:"addr"`
	ENV  string `json:"env"`
}

var config TPV2Config

func main() {

	logger := Config.NewLoggerWithName("main")

	if filePtr, err := os.Open("config/config.json"); err != nil {
		fmt.Printf("文件读取失败:%s,启用默认设置\n", err.Error())
		configStr, _ := Config.TPV2ConfigFs.ReadFile("config.json")
		json.Unmarshal(configStr, &config)
	} else {
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&config)
		filePtr.Close()
	}

	if config.ENV == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	} else if config.ENV == "production" {
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	} else if config.ENV == "trace" {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		panic("未知的环境变量设置")
	}
	logger.Infof("目前环境为：%s", config.ENV)
	r := gin.Default()
	logger.Info("加载中间件")

	for _, v := range Register.MiddlewareCore {
		r.Use(*v)
	}
	logger.Info("中间件加载完成")

	r.Any("/*url", Frame.TinyRouteHandler)
	logger.Infof("TheresaProxyV2启动，地址:%v", config.Addr)

	if err := r.Run(config.Addr); err != nil {
		logger.Panicf("TheresaProxyV2启动失败:%v", err)
	}

}
