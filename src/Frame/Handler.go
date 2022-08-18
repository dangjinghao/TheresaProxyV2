package Frame

import (
	"TheresaProxyV2/src/Library"
	_ "embed"
	"github.com/gin-gonic/gin"
	isDomain "github.com/jbenet/go-is-domain"
	"net/http"
	"strings"
)

// embed
var indexPage string

func TinyRouteHandler(c *gin.Context) {

	requestURI := c.Request.RequestURI
	if requestURI == "/proxy_home" || requestURI == "/proxy_home/" {
		c.String(http.StatusOK, "home")
	} else if strings.HasPrefix(requestURI, "/~/") {
		domainEndIndex := strings.Index(requestURI[3:], "/")
		var requestDomain string
		if domainEndIndex < 0 {
			requestDomain = requestURI[1:]
		} else {
			requestDomain = requestURI[1 : domainEndIndex+1]
		}
		if isDomain.IsDomain(requestDomain) {
			DirectProxyRouter(requestDomain)(c)
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
			Library.ParamProxy(requestDomain)(c)
		} else {
			Library.CookieProxy(c)
		}

	}

}
