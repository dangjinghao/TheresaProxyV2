package plugins

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

type github struct {
	allowedContentTypeSlice []string
	assetGithubReplace      string
	byteGithubReplace       []byte
	byteApiReplace          []byte
}
type ConfigStruct struct {
	ProxySiteScheme string `json:"proxy_site_scheme"`
	ProxySiteDomain string `json:"proxy_site_domain"`
}

func init() {
	var plugin github
	filePtr, err := os.Open("config/github.json")

	if err != nil {
		panic("文件读取失败:" + err.Error())
		return
	}
	var config ConfigStruct
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&config)
	if err != nil {
		panic("decode配置失败：" + err.Error())
		return
	} else {
		plugin.byteGithubReplace = []byte(fmt.Sprintf("%s://%s", config.ProxySiteScheme, config.ProxySiteDomain))
		plugin.byteApiReplace = []byte(fmt.Sprintf("%s://%s/~/api.github.com", config.ProxySiteScheme, config.ProxySiteDomain))
		plugin.assetGithubReplace = fmt.Sprintf("%s://%s/~/objects.githubusercontent.com", config.ProxySiteScheme, config.ProxySiteDomain)

	}
	Register.AddMiddlewareFunc(plugin.RedirectGitClientMiddleware())
	plugin.allowedContentTypeSlice = []string{"html"}
	proxySite := Register.NewProxySiteInfo()
	proxySite.Scheme = "https"
	proxySite.ResponseModify = plugin.ModifyResponse()
	proxySite.AutoGzip = true
	Register.AddProxySite("github.com", proxySite)
	Register.AddProxySite("api.github.com", proxySite)
	Register.AddProxySite("objects.githubusercontent.com", proxySite)
}

func (p *github) ModifyRequest() func(req *http.Request) (err error) {
	return func(req *http.Request) (err error) {
		req.Header.Set("User-Agent", req.Header.Get("User-Agent")+" theresa proxy v2.0.0a1")
		return
	}
}

func (p *github) RedirectGitClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Index(c.Request.Header.Get("User-Agent"), "git") >= 0 &&
			strings.Index(c.Request.URL.Path, "github.com") != 1 {
			//c.String(http.StatusBadRequest, fmt.Sprintf(`git客户端请将URL修改为 "/github.com%v" 而不是 "%v" `, c.Request.URL.Path, c.Request.RequestURI))
			c.Redirect(http.StatusSeeOther, fmt.Sprintf(`/github.com%v`, c.Request.RequestURI))
		} else {
			c.Next()
		}
		return
	}
}
func (p *github) ModifyResponse() func(res *http.Response) (err error) {
	return func(res *http.Response) (err error) {

		if !p.isResponseModifyType(res) {
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

		b = bytes.Replace(b, []byte("https://github.com"), p.byteGithubReplace, -1)
		b = bytes.Replace(b, []byte("https://api.github.com"), p.byteApiReplace, -1)
		res.Header.Set("Location", strings.Replace(res.Header.Get("Location"), "https://objects.githubusercontent.com", p.assetGithubReplace, -1))
		res.Body = io.NopCloser(bytes.NewReader(b))

		return nil
	}
}

func (p *github) isResponseModifyType(res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")
	for _, allowedContentType := range p.allowedContentTypeSlice {
		if strings.Index(contentType, allowedContentType) >= 0 {
			return true
		}
	}
	return false

}
