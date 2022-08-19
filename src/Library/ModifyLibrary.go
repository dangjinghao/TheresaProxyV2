package Library

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func modifyResponseMain(proxyTargetUrl *url.URL) func(res *http.Response) (err error) {

	return func(res *http.Response) (err error) {

		if res.StatusCode >= 400 && res.StatusCode <= 600 {
			//400错误无法修改body会报错
			return nil
		}
		var bodyReader io.ReadCloser
		if res.Header.Get("Content-Encoding") == "gzip" {
			bodyReader, err = gzip.NewReader(res.Body)
			if err != nil {
				if string(err.Error()) == "EOF" {
					return nil
				} else {
					return err
				}
			}
			//当err为EOF时bodyReader并未打开，会导致空指针异常
			defer bodyReader.Close()

			unGzippedBody, err := io.ReadAll(bodyReader)
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
