package Frame

import (
	"TheresaProxyV2/src/Register"
	_ "embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// TODO embed
var indexPage string

func TinyRouteHandler(c *gin.Context) {
	//logger := Config.NewLoggerWithName("tinyRouter")
	requestURI := c.Request.RequestURI

	if requestURI == "/proxy_home" || requestURI == "/proxy_home/" {
		c.String(http.StatusOK, "home")

	} else if pluginRoutePath := c.Request.URL.Path; Register.PluginRoute[pluginRoutePath] != nil {
		//拓展路由
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
		//是否存在
		if Register.ProxySiteCore[requestDomain] != nil || Register.NickNameMap[requestDomain] != "" {

			DirectProxyRouter(requestDomain)(c)
		} else {
			c.String(http.StatusBadRequest, "不允许访问的域名")

		}

		return
	} else {
		//example.com或者非指定域名
		domainEndIndex := strings.Index(requestURI[1:], "/")

		var requestDomain string
		if domainEndIndex < 0 {
			requestDomain = requestURI[1:]
		} else {
			requestDomain = strings.ToLower(requestURI[1 : domainEndIndex+1])
		}
		//是否存在
		if Register.ProxySiteCore[requestDomain] != nil || Register.NickNameMap[requestDomain] != "" {
			ParamProxyRouter(requestDomain)(c)
		} else {
			SessionProxyRouter(requestDomain)(c)

		}

	}

}
