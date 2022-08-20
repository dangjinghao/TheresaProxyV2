package Library

import (
	"TheresaProxyV2/src/Config"
	"TheresaProxyV2/src/Register"
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

// SessionProxy 通过使用cookie的domain值来确定请求域名
func SessionProxy(c *gin.Context) {
	session := sessions.Default(c)
	var proxyDomain string
	if session.Get("domain") == nil {
		c.String(http.StatusBadRequest, "不允许访问")
		return
	} else {
		proxyDomain = session.Get("domain").(string)
	}

	if Register.ProxySiteCore[proxyDomain] == nil {
		c.String(http.StatusBadRequest, "不允许访问的域名")
		return
	}

	proxyTargetUrl := c.Request.URL
	proxyTargetUrl.Host = proxyDomain

	proxyTargetUrl.Scheme = Register.ProxySiteCore[proxyDomain].Scheme
	director := func(req *http.Request) {
		req.Header = c.Request.Header
		req.URL = proxyTargetUrl
		req.Host = req.URL.Host
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ErrorHandler:   noContextCancelErrors,
		ModifyResponse: modifyResponseMain(proxyTargetUrl),
	}
	proxy.ServeHTTP(c.Writer, c.Request)

}

// ParamProxy 通过param参数来获取domain
func ParamProxy(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {

		proxyTargetUrl := c.Request.URL
		proxyTargetUrl.Path = proxyTargetUrl.Path[1+len(proxyDomain):]
		proxyTargetUrl.Host = proxyDomain
		proxyTargetUrl.Scheme = Register.ProxySiteCore[proxyDomain].Scheme
		director := func(req *http.Request) {
			req.Header = c.Request.Header
			req.URL = proxyTargetUrl
			req.Host = req.URL.Host
		}

		proxy := &httputil.ReverseProxy{
			Director:       director,
			ErrorHandler:   noContextCancelErrors,
			ModifyResponse: modifyResponseMain(proxyTargetUrl),
		}

		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

// https://github.com/golang/go/issues/20071#issuecomment-926644055
func noContextCancelErrors(rw http.ResponseWriter, req *http.Request, err error) {
	logger := Config.NewLoggerWithName("ReverseProxyHandler")
	if err != context.Canceled {
		logger.Errorf("http: proxy error: %v", err)
	}
	rw.WriteHeader(http.StatusBadGateway)
}
