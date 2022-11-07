/*
最小可运行例子
*/

package examples

import (
	"TheresaProxyV2/register"
	"net/http"
)

type easiest struct {
}

func (p easiest) Request(req *http.Request) error {
	return nil
}
func (p easiest) Response(res *http.Response) error {
	return nil
}
func init() {
	var p easiest
	site := &register.SiteProperty{}
	site.SiteBehavior = p
	register.ProxySite("httpbin.org", site)

}
