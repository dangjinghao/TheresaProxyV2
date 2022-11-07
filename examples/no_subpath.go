/*
挂载路由/home
*/

package examples

import (
	"TheresaProxyV2/register"
	"github.com/gin-gonic/gin"
	"net/http"
)

type subpath struct {
}

func (p subpath) Request(req *http.Request) error {
	return nil
}
func (p subpath) Response(res *http.Response) error {
	return nil
}
func init() {
	var p subpath
	logger := register.PluginLogger("noSubPath")
	site := &register.SiteProperty{}
	site.SiteBehavior = p
	site.NoDirect = true

	register.Router("/access", func(ctx *gin.Context) {
		err := register.SetTargetDomain(ctx, "example.com")
		if err != nil {
			logger.Errorf("error:%s", err)
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, "/")

	})
	register.ProxySite("example.com", site)

}
