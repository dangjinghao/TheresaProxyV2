package libraries

import (
	"TheresaProxyV2/core"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func responseEditor(targetUrl *url.URL) func(res *http.Response) error {
	return func(res *http.Response) error {
		logger := core.ComponentLogger("modifyResponseMain")
		if res.StatusCode >= 400 && res.StatusCode <= 600 {
			logger.Debugf("请求%s取得异常状态码:%d", targetUrl, res.StatusCode)
			return nil
		}
		var compressCodec codec
		switch res.Header.Get("Content-Encoding") {
		case "gzip":
			logger.Debugf("%s 响应使用gzip,尝试压缩", res.Request.RequestURI)
			compressCodec = gzipCodec{}
		default:
			logger.Debugf("%s 相应未启用或者无法识别的压缩算法，跳过压缩", res.Request.RequestURI)
			compressCodec = defaultCodec{}

		}
		decodedBody, err := compressCodec.decompress(res.Body)
		res.Body = io.NopCloser(decodedBody)
		if err != nil {
			return err
		}
		err = core.ProxySites[targetUrl.Host].Response(res)
		if err != nil {
			return err
		}
		encodedBody, err, length := compressCodec.compress(res.Body)
		res.Body = io.NopCloser(encodedBody)
		if length < 0 {
			logger.Debugf("可能由于未启用解压缩,resbody长度为：%d", length)
			return nil
		}
		res.ContentLength = length
		if res.Header.Get("Content-Length") != "" {
			res.Header.Set("Content-Length", strconv.FormatInt(length, 10))
		}
		return nil

	}
}

type codec interface {
	compress(reader io.Reader) (newReader io.Reader, err error, length int64)
	decompress(reader io.Reader) (newReader io.Reader, err error)
}

type defaultCodec struct {
}

func (m defaultCodec) compress(reader io.Reader) (newReader io.Reader, err error, length int64) {

	byteBody, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	//默认还要重复读取reader一次来确定长度
	length = int64(len(byteBody))
	newReader = io.NopCloser(bytes.NewReader(byteBody))
	return
}
func (m defaultCodec) decompress(reader io.Reader) (newReader io.Reader, err error) {
	return reader, err
}

type gzipCodec struct {
}

func (m gzipCodec) compress(reader io.Reader) (newReader io.Reader, err error, length int64) {
	unGzippedBody, err := io.ReadAll(reader)
	if err != nil {
		return
	}

	var gzipBuffer bytes.Buffer
	gzipBufferWriter := gzip.NewWriter(&gzipBuffer)
	_, err = gzipBufferWriter.Write(unGzippedBody)
	if err != nil {
		return
	}
	if err = gzipBufferWriter.Close(); err != nil {
		return
	}
	newReader = io.NopCloser(bytes.NewReader(gzipBuffer.Bytes()))
	length = int64(gzipBuffer.Len())
	return

}

/*
 */
func (m gzipCodec) decompress(reader io.Reader) (newReader io.Reader, err error) {
	gReader, err := gzip.NewReader(reader)

	if err != nil {
		return nil, err
	}

	defer gReader.Close()
	byteBody, err := io.ReadAll(gReader)

	if err != nil {
		return nil, err
	}

	newReader = io.NopCloser(bytes.NewReader(byteBody))

	return

}
