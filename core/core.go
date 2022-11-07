/*
程序运行核心文件，存储核心结构
*/

package core

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	BindAddr string //框架绑定的地址
	Env      string //运行环境
)

var BaseLogger *logrus.Logger //基logger

type SiteBehavior interface {
	Response(res *http.Response) error
	Request(req *http.Request) error
}

type SiteProperty struct {
	Scheme       string
	AutoCompress bool
	SiteBehavior
}

var (
	ProxySites  = make(map[string]*SiteProperty, 0)
	Nicknames   = make(map[string]string, 0)
	Middlewares = make([]*gin.HandlerFunc, 0)
	Routers     = make(map[string]*gin.HandlerFunc, 0) //储存站点别名到站点域名的映射
	BannedSites = make([]string, 0)                    //禁止子目录方式的反代列表

)

const (
	PathPrefix       = "/~/" //直连反代前缀，理论可替换
	PathPrefixLength = len(PathPrefix)
)

func init() {
	flag.StringVar(&BindAddr, "p", "127.0.0.1:8080", "绑定地址与端口号")
	flag.StringVar(&Env, "e", "dev", "运行环境:trace/dev/prod")
	flag.Parse()

	BaseLogger = logrus.New()
	BaseLogger.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	switch Env {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "dev":
		logrus.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)

	case "prod":
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	default:
		logrus.SetLevel(logrus.DebugLevel)
		Env = "dev"
	}

}
