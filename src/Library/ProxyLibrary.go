package Library

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

// 通过使用cookie的domain值来确定请求域名
func CookieProxy(w http.ResponseWriter, r *http.Request) (err error) {
	cookieProxyDomain, err := r.Cookie("proxy-domain")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "错误的cookie")
		return nil
	}

	cookieDomainStr := cookieProxyDomain.Value
	if Register.ProxySiteCore[cookieDomainStr] == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "不允许访问的域名")
		return nil
	}
	proxyTargetUrl := r.URL
	proxyTargetUrl.Host = cookieDomainStr
	proxyTargetUrl.Scheme = Register.ProxySiteCore[cookieDomainStr].Scheme

	director := func(req *http.Request) {
		req.Header = r.Header
		req.URL = proxyTargetUrl
		req.Host = req.URL.Host
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponseMain(proxyTargetUrl),
	}
	proxy.ServeHTTP(w, r)
	return nil
}

// 通过param参数来获取domain
func ParamProxy(proxyDomain string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxyTargetUrl := r.URL
		proxyTargetUrl.Path = proxyTargetUrl.Path[1+len(proxyDomain):]
		proxyTargetUrl.Host = proxyDomain
		proxyTargetUrl.Scheme = Register.ProxySiteCore[proxyDomain].Scheme

		director := func(req *http.Request) {
			req.Header = r.Header
			req.URL = proxyTargetUrl
			req.Host = req.URL.Host
		}

		proxy := &httputil.ReverseProxy{
			Director:       director,
			ModifyResponse: modifyResponseMain(proxyTargetUrl),
		}
		domainCookie := http.Cookie{
			HttpOnly: true,
			Name:     "proxy-domain",
			Path:     "/",
			Value:    proxyDomain,
			MaxAge:   3600,

			Expires: time.Now().Add(time.Hour),
		}
		http.SetCookie(w, &domainCookie)
		proxy.ServeHTTP(w, r)

	}
}

func modifyResponseMain(proxyTargetUrl *url.URL) func(res *http.Response) (err error) {

	return func(res *http.Response) (err error) {
		//	if r.StatusCode >= 400 && r.StatusCode <= 600 {
		//		//400错误无法修改body会报错
		//		return nil
		//	}
		var bodyReader io.ReadCloser
		if res.Header.Get("Content-Encoding") == "gzip" {
			bodyReader, err = gzip.NewReader(res.Body)
			defer bodyReader.Close()

			if err != nil {
				return err
			}
			unGzippedBody, err := io.ReadAll(bodyReader)
			if err != nil {
				fmt.Println(err)
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
			if res.Header.Get("Content-Length") != "" {
				res.Header.Set("Content-Length", strconv.Itoa(gzipBuffer.Len()))
			}
			err = res.Body.Close()
			if err != nil {
				return err
			}
		} else {
			//未启用gzip
			err = Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify(res)
		}

		return nil
	}

}
