package Library

import (
	"TheresaProxyV2/src/Register"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

// 通过使用cookie的domain值来确定请求域名
func CookieProxy(c *gin.Context) (err error) {
	cookieProxyDomain, err := c.Cookie("proxy-domain")
	if err != nil {
		c.String(http.StatusBadRequest, "错误的cookie")
		return nil
	}
	if Register.ProxySiteCore[cookieProxyDomain] == nil {
		c.String(http.StatusBadRequest, "不允许访问的域名")
		return nil
	}
	if Register.MiddlewareCore[cookieProxyDomain] != nil {
		(*Register.MiddlewareCore[cookieProxyDomain])(c)
	}
	proxyTargetUrl := c.Request.URL
	proxyTargetUrl.Host = cookieProxyDomain

	proxyTargetUrl.Scheme = Register.ProxySiteCore[cookieProxyDomain].Scheme
	if Register.ProxySiteCore[cookieProxyDomain].RequestModify != nil {
		Register.ProxySiteCore[cookieProxyDomain].RequestModify(c.Request)
	}

	director := func(req *http.Request) {
		req.Header = c.Request.Header
		req.URL = proxyTargetUrl
		req.Host = req.URL.Host
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
		ErrorHandler: func(http.ResponseWriter, *http.Request, error) {
			// handle error
		},
		ModifyResponse: modifyResponseMain(proxyTargetUrl),
	}
	proxy.ServeHTTP(c.Writer, c.Request)
	return nil
}

// 通过param参数来获取domain
func ParamProxy(proxyDomain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if Register.MiddlewareCore[proxyDomain] != nil {
			(*Register.MiddlewareCore[proxyDomain])(c)
		}

		proxyTargetUrl := c.Request.URL
		proxyTargetUrl.Path = proxyTargetUrl.Path[1+len(proxyDomain):]
		proxyTargetUrl.Host = proxyDomain
		proxyTargetUrl.Scheme = Register.ProxySiteCore[proxyDomain].Scheme
		if Register.ProxySiteCore[proxyDomain].RequestModify != nil {
			Register.ProxySiteCore[proxyDomain].RequestModify(c.Request)
		}

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
		_, err := c.Cookie("proxy-domain")
		if err != nil {
			c.SetCookie("proxy-domain", proxyDomain, 3600, "/", "", true, false)
		}
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

// https://github.com/golang/go/issues/20071#issuecomment-926644055
func noContextCancelErrors(rw http.ResponseWriter, req *http.Request, err error) {
	if err != context.Canceled {

		fmt.Errorf("http: proxy error: %v", err)
	}
	rw.WriteHeader(http.StatusBadGateway)
}
