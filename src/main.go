package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Frame"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Any("/*url", Frame.TinyRouteHandler)
	r.Run()

}
