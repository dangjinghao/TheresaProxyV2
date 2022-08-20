package Library

import (
	"TheresaProxyV2/src/Config"
	"TheresaProxyV2/src/Register"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func modifyResponseMain(proxyTargetUrl *url.URL) func(res *http.Response) (err error) {
	logger := Config.NewLoggerWithName("modifyResponseMain")
	return func(res *http.Response) (err error) {

		if res.StatusCode >= 400 && res.StatusCode <= 600 {
			logger.Infof("请求%q取得异常状态码:%v", res.Request.RequestURI, res.StatusCode)
			return nil
		}
		var bodyReader io.ReadCloser
		if res.Header.Get("Content-Encoding") == "gzip" {
			logger.Debugf("%s 响应使用gzip,尝试压缩", res.Request.RequestURI)
			bodyReader, err = gzip.NewReader(res.Body)
			if err != nil {
				if string(err.Error()) == "EOF" {
					logger.Debugf("%q响应体为EOF", res.Request.RequestURI)

					return nil
				} else {
					return err
				}
			}
			//当err为EOF时bodyReader并未打开，会导致空指针异常
			defer bodyReader.Close()

			unGzippedBody, err := io.ReadAll(bodyReader)
			if err != nil {
				logger.Debugf("%q响应体ungzip失败，异常为%v", res.Request.RequestURI, err.Error())
				return err
			}

			res.Body = io.NopCloser(bytes.NewReader(unGzippedBody))
			if Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify != nil {
				err = Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify(res)
				if err != nil {
					logger.Debugf("%q调用对应域名responseModify失败，异常为%v", res.Request.RequestURI, err.Error())
					return err
				}
			}

			//重新生成body以启用gzip压缩
			unGzippedBody, err = io.ReadAll(res.Body)
			if err != nil {
				logger.Debugf("%q再读取ungzip的body失败，异常为%v", res.Request.RequestURI, err.Error())
				return err
			}
			var gzipBuffer bytes.Buffer
			w := gzip.NewWriter(&gzipBuffer)
			_, err = w.Write(unGzippedBody)
			if err != nil {
				logger.Debugf("向gzipbuffer中写入gzip压缩后的的body失败,异常为%v", err.Error())

				return err
			}

			err = w.Close()

			res.Body = io.NopCloser(bytes.NewReader(gzipBuffer.Bytes()))
			res.ContentLength = int64(gzipBuffer.Len())
			if res.Header.Get("Content-Length") != "" {
				res.Header.Set("Content-Length", strconv.Itoa(gzipBuffer.Len()))
			}

			//确切说我不知道到底应不应该关闭这个
			err = res.Body.Close()
			if err != nil {
				logger.Debugf("关闭%v的body失败,异常为%v", res.Request.RequestURI, err.Error())
				return err
			}
		} else {
			logger.Debugf("%q响应未使用gzip,跳过压缩", res.Request.RequestURI)
			if Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify != nil {
				err = Register.ProxySiteCore[proxyTargetUrl.Host].ResponseModify(res)
				if err != nil {
					logger.Debugf("%q调用对应域名responseModify失败，异常为%v", res.Request.RequestURI, err.Error())
					return err
				}
			}
		}

		return nil
	}

}
