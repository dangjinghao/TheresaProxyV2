package plugins

import (
	"TheresaProxyV2/register"
	"github.com/sirupsen/logrus"
	"net/http"
)

type apocalypse struct {
	logger *logrus.Entry
}

func (p apocalypse) Request(*http.Request) error {
	return nil
}
func (p apocalypse) Response(*http.Response) error {
	return nil
}
func init() {
	var p apocalypse
	p.logger = register.PluginLogger("apocalypse")
	p.logger.Info("teri teri~")
}
