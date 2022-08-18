package main

import (
	_ "TheresaProxyV2/plugins"
	"TheresaProxyV2/src/Library"
	"TheresaProxyV2/src/Route"
	"fmt"
	isDomain "github.com/jbenet/go-is-domain"
	"net/http"
	"strings"
)

func tinyRouteHandler(w http.ResponseWriter, r *http.Request) {

	requestURI := r.RequestURI
	if requestURI == "/home" || requestURI == "/home/" {
		fmt.Fprintf(w, "home")
	} else if strings.HasPrefix(requestURI, "/~/") {
		domainEndIndex := strings.Index(requestURI[3:], "/")
		var requestDomain string
		if domainEndIndex < 0 {
			requestDomain = requestURI[1:]
		} else {
			requestDomain = requestURI[1 : domainEndIndex+1]
		}
		if isDomain.IsDomain(requestDomain) {
			Route.DirectProxyRouter(requestDomain)(w, r)
		}
		return
	} else {
		domainEndIndex := strings.Index(requestURI[1:], "/")
		var requestDomain string
		if domainEndIndex < 0 {
			requestDomain = requestURI[1:]
		} else {
			requestDomain = requestURI[1 : domainEndIndex+1]
		}
		if isDomain.IsDomain(requestDomain) {
			Library.ParamProxy(requestDomain)(w, r)
		} else {
			Library.CookieProxy(w, r)
		}

	}

}

func main() {

	http.HandleFunc("/", tinyRouteHandler)
	http.ListenAndServe(":8080", nil)
}
