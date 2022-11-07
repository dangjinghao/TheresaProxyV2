/*
挂载路由/home
*/

package examples

import (
	"TheresaProxyV2/register"
	"github.com/gin-gonic/gin"
	"net/http"
)

type home struct {
}

func (p home) Request(req *http.Request) error {
	return nil
}
func (p home) Response(res *http.Response) error {
	return nil
}
func init() {
	var p home
	site := &register.SiteProperty{}
	site.SiteBehavior = p
	register.Router("/home", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})
	register.ProxySite("httpbin.org", site)

}
