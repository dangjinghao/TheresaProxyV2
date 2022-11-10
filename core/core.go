/*
程序运行核心文件，存储核心结构
*/

package core

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	BindAddr     string                       //框架绑定的地址
	Env          string                       //运行环境
	additionFlag string                       //拓展flag
	FlagDict     = make(map[string]string, 0) //用于插件的拓展flag字典
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

func dumpAdditionFlags(additionFlag string, flagDict map[string]string) {
	flagArray := strings.Split(additionFlag, ";")
	for _, i := range flagArray {
		if i != "" {
			_itemSlice := strings.Split(i, "=")
			flagDict[_itemSlice[0]] = _itemSlice[1]
		}

	}
}

func init() {
	flag.StringVar(&BindAddr, "p", "127.0.0.1:8080", "绑定地址与端口号")
	flag.StringVar(&Env, "e", "dev", "运行环境:trace/dev/prod")
	flag.StringVar(&additionFlag, "add", "", "可以由插件识别的拓展flag，"+
		"格式:[name1]=[value1];[name2]=[value2];...")
	flag.Parse()

	BaseLogger = logrus.New()
	BaseLogger.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	dumpAdditionFlags(additionFlag, FlagDict)
	switch Env {
	case "trace":
		BaseLogger.SetLevel(logrus.TraceLevel)
	case "dev":
		BaseLogger.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	case "prod":
		BaseLogger.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	default:
		BaseLogger.SetLevel(logrus.DebugLevel)
		Env = "dev"
		gin.SetMode(gin.DebugMode)
	}

}
