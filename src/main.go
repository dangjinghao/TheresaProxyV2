package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Route"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	r := gin.Default()
	r.Any("/:domain/*proxyPath", Route.ProxyCheck)
	r.Any("/", Route.RootProxy)

	r.GET("/home", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.Any("/~/:domain/*proxyPath", Route.DirectProxy)
	r.Run()
}
