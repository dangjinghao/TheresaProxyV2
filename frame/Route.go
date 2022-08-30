package frame

import (
	"TheresaProxyV2/library"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func DirectProxyRouter(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Request.URL.Path = c.Request.URL.Path[2:]
		library.ParamProxy(proxyDomain)(c)
	}
}
func SessionProxyRouter(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		library.SessionProxy(c)
	}
}

func ParamProxyRouter(requestDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("domain") == nil || session.Get("domain").(string) != requestDomain {
			session.Set("domain", requestDomain)
			session.Save()
		}
		library.ParamProxy(requestDomain)(c)

	}
}
