/*
挂载中间件示例
*/
package examples

import (
	"TheresaProxyV2/register"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type pong struct {
	logger *logrus.Entry
}

func (p pong) Request(req *http.Request) error {
	p.logger.Info("req pong")
	return nil
}
func (p pong) Response(res *http.Response) error {
	p.logger.Info("res pong")
	return nil
}
func init() {
	var p pong
	site := &register.SiteProperty{}
	p.logger = register.PluginLogger("pong")
	site.SiteBehavior = p
	register.Middleware(func(ctx *gin.Context) {
		p.logger.Infof("pong!")
	})
	register.ProxySite("httpbin.org", site)

}
