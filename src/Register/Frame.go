package Register

import (
	"github.com/gin-gonic/gin"
)

var MiddlewareCore []*gin.HandlerFunc

func AddMiddlewareFunc(middlewareFunc gin.HandlerFunc) {
	MiddlewareCore = append(MiddlewareCore, &middlewareFunc)
	return
}
func init() {
	MiddlewareCore = make([]*gin.HandlerFunc, 0)
}
