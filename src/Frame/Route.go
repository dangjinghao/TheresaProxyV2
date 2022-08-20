package Frame

import (
	"TheresaProxyV2/src/Library"
	"TheresaProxyV2/src/Register"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DirectProxyRouter(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if Register.ProxySiteCore[proxyDomain] == nil {
			c.String(http.StatusBadRequest, "不允许访问的域名")
			return
		}
		Library.ParamProxy(proxyDomain)(c)
	}
}
func SessionProxyRouter(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		Library.SessionProxy(c)
	}
}

func ParamProxyRouter(requestDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("domain") == nil || session.Get("domain").(string) != requestDomain {
			session.Set("domain", requestDomain)
			session.Save()
		}
		Library.ParamProxy(requestDomain)(c)

	}
}
