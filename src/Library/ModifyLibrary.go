package Library

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

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
