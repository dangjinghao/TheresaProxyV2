package libraries

import (
	"TheresaProxyV2/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func tinyRouter(c *gin.Context) {

	requestURI := c.Request.RequestURI

	//处理挂载路由
	if routerPath := c.Request.URL.Path; core.Routers[routerPath] != nil {
		(*core.Routers[routerPath])(c)

	} else if strings.HasPrefix(requestURI, core.PathPrefix) {
		//存在前缀时的不修改session直接反代
		domainEnd := strings.Index(requestURI[core.PathPrefixLength:], "/")

		if domainEnd < 0 {
			// /~/xxx
			c.String(http.StatusBadRequest, "不允许访问")
			return
		}
		//取得无session反代的domain
		domain := routerPath[core.PathPrefixLength : domainEnd+core.PathPrefixLength]

		//处理/~/[domain]
		if core.ExistDomain(domain) {
			SubPathAdapter(domain)(c)
		} else {
			//TODO 不允许访问
			return
		}

	} else { // /[domain]xxx

		var domain string
		if domainEnd := strings.Index(requestURI[1:], "/"); domainEnd < 0 {
			// /[domain]
			domain = requestURI[1:]
		} else {
			// /[domain]/xxx
			domain = strings.ToLower(requestURI[1 : domainEnd+1])
		}
		if core.ExistDomain(domain) {
			SubPathAdapter(domain)(c)
		} else {
			SessionAdapter(c)
		}
	}

}
func InitRouter(bindAddr string) {
	//TODO gin logger
	logger := core.ComponentLogger("InitRouter")
	r := gin.New()
	logger.Debug("加载session中间件")

	store := cookie.NewStore([]byte("proxy-secret"))
	r.Use(sessions.Sessions("proxy-session", store))

	logger.Debug("加载自定义格式gin日志中间件")
	r.Use(ginLogger())
	logger.Debug("加载拓展中间件")
	for _, v := range core.Middlewares {
		r.Use(*v)
	}

	r.Any("/*url", tinyRouter)
	logger.Infof("应用已在地址:%s 启动", bindAddr)

	if err := r.Run(bindAddr); err != nil {
		logger.Panicf("启动失败，Err:%s", err)
	}
}

func ginLogger() gin.HandlerFunc {
	logger := core.ComponentLogger("web-framework")

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime).Microseconds()
		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"Code":   statusCode,
			"Time":   latencyTime,
			"ip":     clientIP,
			"Method": reqMethod,
			"URI":    reqUrl,
		}).Info("done")

	}
}
