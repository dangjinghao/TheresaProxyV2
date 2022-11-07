package main

import (
	"TheresaProxyV2/core"
	"TheresaProxyV2/libraries"
	_ "TheresaProxyV2/plugins"
)

var Version = "unknown" //核心版本，编译时设置

func main() {
	logger := core.ComponentLogger("main")
	logger.Infof("版本:%s", Version)
	logger.Infof("运行模式:%s", core.Env)
	libraries.InitRouter(core.BindAddr)
}
