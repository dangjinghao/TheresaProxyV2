package Route

import (
	"TheresaProxyV2/src/Library"
	"TheresaProxyV2/src/Register"
	"fmt"
	isDomain "github.com/jbenet/go-is-domain"
	"net/http"
)

func DirectProxyRouter(proxyDomain string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Register.ProxySiteCore[proxyDomain] == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "不允许访问的域名")
			return
		}
		if isDomain.IsDomain(proxyDomain) {
			//为url中包含domain且未设定cookie的请求设置domain cookie
			proxyTargetUrl := r.URL

			proxyTargetUrl.Path = proxyTargetUrl.Path[3+len(proxyDomain):]
			Library.ParamProxy(proxyDomain)(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "错误的域名")
			return
		}
	}
}
