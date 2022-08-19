package Register

import (
	"TheresaProxyV2/src/Config"
	"github.com/sirupsen/logrus"
	"net/http"
)

var ProxySiteCore map[string]*ProxySiteInfo
var logger *logrus.Entry

type ProxySiteInfo struct {
	Scheme         string
	AutoGzip       bool
	ResponseModify func(*http.Response) error
	RequestModify  func(*http.Request) error
}

func NewProxySiteInfo() ProxySiteInfo {

	return ProxySiteInfo{}
}
func AddProxySite(domain string, proxySite ProxySiteInfo) {
	logger.Infof("添加代理站点:%v", domain)
	ProxySiteCore[domain] = &proxySite
	return
}
func init() {
	logger = Config.NewLoggerWithName("netRegister")

	ProxySiteCore = make(map[string]*ProxySiteInfo, 0)
}
