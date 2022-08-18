package Route

import (
	"TheresaProxyV2/src/Register"
	"github.com/gin-gonic/gin"
	isDomain "github.com/jbenet/go-is-domain"
	"net/http"
	"strings"
)

func ProxyCheck(c *gin.Context) {

	if siteDomain := c.Param("domain"); isDomain.IsDomain(siteDomain) {
		if Register.ProxySiteCore[siteDomain] == nil {
			c.String(http.StatusBadRequest, "不允许请求的网址")
			return
		}
		//TODO delete the consume domain= cookie
		//TODO middleware hook support
		//TODO request modify support
		//TODO plugin loader and running logger
		//为url中包含domain且未设定cookie的请求设置domain cookie
		c.SetCookie("domain", siteDomain, 3600, "/", "", false, true)
		proxyTargetUrl := c.Request.URL

		ParamProxy(proxyTargetUrl, c)
	} else {
		cookieDomain, getCookieErr := c.Cookie("domain")
		//没有domain cookie就返回
		if getCookieErr != nil {
			c.String(http.StatusBadRequest, "cookie错误")
			return
		}
		if Register.ProxySiteCore[cookieDomain] == nil {
			c.String(http.StatusBadRequest, "不允许请求的网址")
			return
		}
		//在有domain cookie的情况下取domain cookie直接请求
		CookieProxy(c)

	}

}

func DirectProxy(c *gin.Context) {
	siteDomain := c.Param("domain")
	if Register.ProxySiteCore[siteDomain] == nil {
		c.String(http.StatusBadRequest, "不允许请求的网址")
		return
	}
	if isDomain.IsDomain(siteDomain) {
		//为url中包含domain且未设定cookie的请求设置domain cookie
		proxyTargetUrl := c.Request.URL

		proxyTargetUrl.Path = strings.Replace(proxyTargetUrl.Path, "/~", "", 1)
		ParamProxy(proxyTargetUrl, c)
	} else {
		c.String(http.StatusBadRequest, "错误的请求地址")
		return
	}
}

func RootProxy(c *gin.Context) {
	cookieDomain, err := c.Cookie("domain")
	if err != nil {
		c.String(http.StatusBadRequest, "cookie非法")
		return
	}
	if Register.ProxySiteCore[cookieDomain] == nil {
		c.String(http.StatusBadRequest, "不允许请求的网址")
		return
	}

	CookieProxy(c)
}
