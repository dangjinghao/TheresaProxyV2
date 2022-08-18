package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Frame"
	"TheresaProxyV2/src/Register"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	for _, v := range Register.MiddlewareCore {
		r.Use(*v)
	}

	r.Use()
	r.Any("/*url", Frame.TinyRouteHandler)

	r.Run()

}
