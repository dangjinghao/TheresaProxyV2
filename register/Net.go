package register

import (
	"TheresaProxyV2/rawConfig"
	"github.com/sirupsen/logrus"
	"net/http"
)

var ProxySiteCore map[string]*ProxySiteInfo
var NickNameMap map[string]string
var logger *logrus.Entry

type ProxySiteInfo struct {
	Scheme         string
	Nickname       string
	AutoGzip       bool
	ResponseModify func(*http.Response) error
	RequestModify  func(*http.Request) error
}

func NewProxySiteInfo() ProxySiteInfo {

	return ProxySiteInfo{}
}
func AddProxySite(domain string, proxySite ProxySiteInfo) {
	if proxySite.Scheme == "" {
		proxySite.Scheme = "http"
		logger.Errorf("站点%v缺少scheme，将其设置为http", domain)
	}
	logger.Infof("添加代理站点%v", domain)
	ProxySiteCore[domain] = &proxySite
	if proxySite.Nickname != "" {
		logger.Infof("添加代理站点%v别名%v", domain, proxySite.Nickname)
		NickNameMap[proxySite.Nickname] = domain
	}
	return
}
func init() {
	logger = rawConfig.NewLoggerWithName("netRegister")
	ProxySiteCore = make(map[string]*ProxySiteInfo, 0)
	NickNameMap = make(map[string]string, 0)
}
