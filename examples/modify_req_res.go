/*
修改请求和响应
*/

package examples

import (
	"TheresaProxyV2/register"
	"bytes"
	"io"
	"net/http"
)

type httpbin2 struct {
}

func (p httpbin2) Request(req *http.Request) error {
	req.Header.Set("User-Agent", req.Header.Get("User-Agent")+" theresa")
	return nil
}
func (p httpbin2) Response(res *http.Response) error {
	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	buffer = bytes.Replace(buffer, []byte("args"), []byte("argvs"), -1)
	res.Body = io.NopCloser(bytes.NewReader(buffer))
	return nil
}
func init() {
	var p httpbin2
	site := &register.SiteProperty{}
	site.SiteBehavior = p
	register.ProxySite("httpbin.org", site)

}
