package Register

import (
	"github.com/gin-gonic/gin"
)

var MiddlewareCore map[string]*gin.HandlerFunc

func AddMiddlewareFunc(domain string, middlewareFunc gin.HandlerFunc) {
	MiddlewareCore[domain] = &middlewareFunc
	return
}
func init() {
	MiddlewareCore = make(map[string]*gin.HandlerFunc, 0)
}
