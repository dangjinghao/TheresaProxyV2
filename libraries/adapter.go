package libraries

import (
	"TheresaProxyV2/core"
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"

	"net/http"
	"net/http/httputil"
)

func SessionAdapter(c *gin.Context) {
	session := sessions.Default(c)
	var domain string
	if session.Get("domain") == nil {
		c.String(http.StatusBadRequest, "不允许访问")

		return
	} else {
		domain = session.Get("domain").(string)
	}
	if core.Nicknames[domain] != "" {
		domain = core.Nicknames[domain]
	}
	if core.ProxySites[domain] == nil {
		c.String(http.StatusBadRequest, "不允许访问")
	}
	targetUrl := c.Request.URL
	targetUrl.Host = domain
	targetUrl.Scheme = core.ProxySites[domain].Scheme
	director := func(req *http.Request) {
		req.Header = c.Request.Header
		req.URL = targetUrl
		req.Host = req.URL.Host
	}
	err := core.ProxySites[domain].Request(c.Request)
	if err != nil {
		return
	}
	proxy := &httputil.ReverseProxy{
		Director:       director,
		ErrorHandler:   noCtxErr,
		ModifyResponse: responseEditor(targetUrl),
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func SubPathAdapter(domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("domain") == nil || session.Get("domain").(string) != domain {
			session.Set("domain", domain)
			err := session.Save()
			if err != nil {
				return
			}
		}

		targetUrl := c.Request.URL
		//如果包含PathPrefix，将其删除
		if strings.HasPrefix(targetUrl.Path, core.PathPrefix) {
			targetUrl.Path = targetUrl.Path[1+len(domain):]
		}
		//删除域名子路径
		targetUrl.Path = targetUrl.Path[1+len(domain):]
		if core.Nicknames[domain] != "" {
			domain = core.Nicknames[domain]
		}
		if core.InSlice[string](domain, core.BannedSites) {
			c.String(http.StatusBadRequest, "不允许访问")
		}
		targetUrl.Host = domain
		targetUrl.Scheme = core.ProxySites[domain].Scheme
		director := func(req *http.Request) {
			req.Header = c.Request.Header
			req.URL = targetUrl
			req.Host = req.URL.Host
		}
		err := core.ProxySites[domain].Request(c.Request)
		if err != nil {
			return
		}
		proxy := &httputil.ReverseProxy{
			Director:       director,
			ErrorHandler:   noCtxErr,
			ModifyResponse: responseEditor(targetUrl),
		}
		proxy.ServeHTTP(c.Writer, c.Request)

	}
}

// https://github.com/golang/go/issues/20071#issuecomment-926644055
func noCtxErr(rw http.ResponseWriter, req *http.Request, err error) {
	logger := core.ComponentLogger("ReverseProxyHandler")
	if err != context.Canceled {
		logger.Errorf("反代错误 %s", err)
	}
	rw.WriteHeader(http.StatusBadGateway)
}
