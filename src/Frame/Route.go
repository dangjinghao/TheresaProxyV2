package Frame

import (
	"TheresaProxyV2/src/Library"
	"TheresaProxyV2/src/Register"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DirectProxyRouter(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if Register.ProxySiteCore[proxyDomain] == nil {
			c.String(http.StatusBadRequest, "不允许访问的域名")
			return
		}
		//为url中包含domain且未设定cookie的请求设置domain cookie
		proxyTargetUrl := c.Request.URL
		proxyTargetUrl.Path = proxyTargetUrl.Path[2:]

		Library.ParamProxy(proxyDomain)(c)
	}
}
