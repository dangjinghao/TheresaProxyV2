/*
开启AutoCompress后，将会自动对gzip网页解码编码以修改响应
*/

package examples

import (
	"TheresaProxyV2/register"
	"bytes"
	"io"
	"net/http"
)

type gzipCompress struct {
}

func (p gzipCompress) Request(req *http.Request) error {
	return nil
}
func (p gzipCompress) Response(res *http.Response) error {
	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	buffer = bytes.Replace(buffer, []byte("Intro"), []byte("Introduction"), 1)
	res.Body = io.NopCloser(bytes.NewReader(buffer))
	return nil
}
func init() {
	var p gzipCompress
	site := &register.SiteProperty{
		Scheme:       "https",
		Nickname:     "gzip",
		AutoCompress: true,
	}
	site.SiteBehavior = p
	register.ProxySite("www.gzip.org", site)
}
