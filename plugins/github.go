package plugins

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"io"
	"net/http"
	"strings"
)

type github struct {
	allowedContentTypeSlice []string
}

func init() {
	var plugin github
	plugin.allowedContentTypeSlice = []string{"html"}
	proxySite := Register.NewProxySiteInfo()
	proxySite.Scheme = "https"

	proxySite.ResponseModify = plugin.ModifyResponse()
	proxySite.AutoGzip = true
	Register.AddProxySite("github.com", proxySite)
	Register.AddProxySite("api.github.com", proxySite)
}

func (p *github) ModifyResponse() func(res *http.Response) (err error) {
	return func(res *http.Response) (err error) {

		if !p.isResponseModified(res) {
			return nil
		}
		delete(res.Header, "Content-Security-Policy")
		//res.Request.Header.Set("Referer", "https://github.com")
		bodyReader := res.Body
		b, err := io.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		err = res.Body.Close()
		if err != nil {
			return err
		}

		b = bytes.Replace(b, []byte("https://github.com"), []byte("http://127.0.0.1:8080"), -1)
		b = bytes.Replace(b, []byte("https://api.github.com"), []byte("http://127.0.0.1:8080/~/api.github.com"), -1)

		res.Body = io.NopCloser(bytes.NewReader(b))

		return nil
	}
}

func (p *github) isResponseModified(res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")
	for _, allowedContentType := range p.allowedContentTypeSlice {
		if strings.Index(contentType, allowedContentType) >= 0 {
			return true
		}
	}
	return false

}
