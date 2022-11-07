package examples

import (
	"TheresaProxyV2/register"
	"net/http"
)

type fullFunction struct {
}

func (p fullFunction) Request(req *http.Request) error {
	return nil
}
func (p fullFunction) Response(res *http.Response) error {
	return nil
}
func init() {
	var p fullFunction
	site := &register.SiteProperty{
		Scheme:       "https",
		Nickname:     "hb",
		AutoCompress: true,
	}
	site.SiteBehavior = p
	register.ProxySite("httpbin.org", site)
}
