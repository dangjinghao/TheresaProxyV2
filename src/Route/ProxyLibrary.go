package Route

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

//通过使用cookie的domain值来确定请求域名
func CookieProxy(c *gin.Context) {
	cookieDomain, _ := c.Cookie("domain")

	proxyTargetUrl := c.Request.URL
	proxyTargetUrl.Host = cookieDomain
	proxyTargetUrl.Scheme = Register.ProxySiteCore[cookieDomain].Scheme

	director := func(req *http.Request) {
		req.Header = c.Request.Header
		req.URL = proxyTargetUrl
		req.Host = req.URL.Host
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponseMain(proxyTargetUrl),
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

//通过param参数来获取domain
func ParamProxy(proxyTargetUrl *url.URL, c *gin.Context) {
	siteDomain := c.Param("domain")

	proxyTargetUrl.Path = strings.Replace(c.Request.URL.Path, "/"+siteDomain, "", 1)
	proxyTargetUrl.Host = siteDomain
	proxyTargetUrl.Scheme = Register.ProxySiteCore[siteDomain].Scheme

	director := func(req *http.Request) {
		req.Header = c.Request.Header
		req.URL = proxyTargetUrl
		req.Host = req.URL.Host
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponseMain(proxyTargetUrl),
	}
	proxy.ServeHTTP(c.Writer, c.Request)

}

func modifyResponseMain(proxyTargetUrl *url.URL) func(res *http.Response) (err error) {

	return func(res *http.Response) (err error) {
		if res.StatusCode >= 400 && res.StatusCode <= 600 {
			//400错误无法修改body会报错
			return nil
		}

		var bodyReader io.ReadCloser
		if res.Header.Get("Content-Encoding") == "gzip" {
			bodyReader, err = gzip.NewReader(res.Body)
			defer bodyReader.Close()

			unGzippedBody, err := io.ReadAll(bodyReader)
			if err != nil {
				return err
			}
			err = res.Body.Close()
			if err != nil {
				return err
			}
			res.Body = io.NopCloser(bytes.NewReader(unGzippedBody))
			err = Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify(res)

			//重新生成body以启用gzip压缩
			unGzippedBody, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			var gzipBuffer bytes.Buffer
			w := gzip.NewWriter(&gzipBuffer)
			w.Write(unGzippedBody)
			w.Close()

			res.Body = io.NopCloser(bytes.NewReader(gzipBuffer.Bytes()))
			res.ContentLength = int64(gzipBuffer.Len())

			res.Header.Set("Content-Length", strconv.Itoa(gzipBuffer.Len()))

		} else {
			//未启用gzip
			err = Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify(res)
		}

		return nil
	}

}
