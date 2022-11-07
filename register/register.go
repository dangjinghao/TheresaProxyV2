package register

import (
	"TheresaProxyV2/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
)

var logger = core.ComponentLogger("register")

type SiteProperty struct {
	Scheme       string //scheme
	Nickname     string //站点别名
	AutoCompress bool   //启用自动解压压缩
	NoDirect     bool   //如果为true，将无法通过访问子目录方式反代
	core.SiteBehavior
}

func PluginLogger(name string) *logrus.Entry {
	return core.BaseLogger.WithFields(logrus.Fields{"type": "plugin", "name": name})
}

func Router(path string, handlerFunc gin.HandlerFunc) {
	if !strings.HasPrefix(path, "/") {
		logger.Errorf("%s路由挂载失败:缺少'/'前缀", path)
		return
	}
	core.Routers[path] = &handlerFunc
	logger.Debugf("成功挂载插件路由:%s", path)
}

func Middleware(handlerFunc gin.HandlerFunc) {
	core.Middlewares = append(core.Middlewares, &handlerFunc)

}

func ProxySite(target string, property *SiteProperty) {
	if property.Scheme == "" {
		property.Scheme = "http"
		logger.Debugf("%s的属性为空，设置为http", target)
	}

	core.ProxySites[target] = &core.SiteProperty{
		Scheme:       property.Scheme,
		AutoCompress: property.AutoCompress,
		SiteBehavior: property.SiteBehavior,
	}
	if property.Nickname != "" {
		core.Nicknames[property.Nickname] = target
		logger.Debugf("%s映射到%s", property.Nickname, target)
	}
	if property.NoDirect {
		if property.Nickname != "" {
			logger.Errorf("%s已存在别名%s但不被允许通过子目录反代，将强制禁止子目录反代", target, property.Nickname)
		}
		core.BannedSites = append(core.BannedSites, target)
		logger.Debugf("已把%s添加入禁止子目录反代列表", target)
	}
}
