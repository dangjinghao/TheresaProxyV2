package Register

import (
	"net/http"
)

var ProxySiteCore map[string]*ProxySiteInfo

type ProxySiteInfo struct {
	Scheme         string
	AutoGzip       bool
	ResponseModify func(*http.Response) error
	RequestModify  func(*http.Request)
}

func NewProxySiteInfo() ProxySiteInfo {

	return ProxySiteInfo{}
}
func AddProxySite(domain string, proxySite ProxySiteInfo) {
	ProxySiteCore[domain] = &proxySite
	return
}
func init() {
	ProxySiteCore = make(map[string]*ProxySiteInfo, 0)
}
