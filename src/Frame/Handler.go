package Frame

import (
	"TheresaProxyV2/src/Config"
	"TheresaProxyV2/src/Register"
	_ "embed"
	"github.com/gin-gonic/gin"
	isDomain "github.com/jbenet/go-is-domain"
	"net/http"
	"strings"
)

// TODO embed
var indexPage string

func TinyRouteHandler(c *gin.Context) {
	logger := Config.NewLoggerWithName("tinyRouter")
	requestURI := c.Request.RequestURI

	if requestURI == "/proxy_home" || requestURI == "/proxy_home/" {
		c.String(http.StatusOK, "home")

	} else if pluginRoutePath := c.Request.URL.Path; Register.PluginRoute[pluginRoutePath] != nil {
		(*Register.PluginRoute[pluginRoutePath])(c)

	} else if strings.HasPrefix(requestURI, "/~/") {
		//直接代理
		domainEndIndex := strings.Index(requestURI[3:], "/")

		var requestDomain string
		if domainEndIndex < 0 {
			requestDomain = requestURI[1:]
		} else {
			requestDomain = requestURI[3 : domainEndIndex+3]
		}

		if isDomain.IsDomain(requestDomain) {
			DirectProxyRouter(requestDomain)(c)
		} else {
			logger.Debugf("错误的域名:%v", requestDomain)
			c.String(http.StatusBadRequest, "不合规请求")
		}

		return
	} else {
		domainEndIndex := strings.Index(requestURI[1:], "/")

		var requestDomain string
		if domainEndIndex < 0 {

			requestDomain = requestURI[1:]
		} else {

			requestDomain = strings.ToLower(requestURI[1 : domainEndIndex+1])
		}
		if isDomain.IsDomain(requestDomain) {
			ParamProxyRouter(requestDomain)(c)
		} else {
			SessionProxyRouter(requestDomain)(c)
		}

	}

}
