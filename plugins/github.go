package plugins

import (
	"TheresaProxyV2/src/Register"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type github struct {
	allowedContentTypeSlice     []string
	stringObjectsContentReplace string
	byteGithubReplace           string
	byteRawReplace              string
	byteApiReplace              string
}

var githubLogger *logrus.Entry

type githubConfig struct {
	ProxySiteScheme string `json:"proxy_site_scheme"`
	ProxySiteDomain string `json:"proxy_site_domain"`
}

func init() {
	var plugin github
	var config githubConfig
	githubLogger = Register.GetPluginLogger("github")
	plugin.getConfig(&config)

	plugin.loadReplaceStr(config)
	plugin.allowedContentTypeSlice = []string{"html"}
	plugin.addToRegister()

}

func (p *github) getConfig(config *githubConfig) {
	filePtr, err := Register.GetPluginConfig("github")

	if err != nil {
		githubLogger.Panic("文件读取失败:" + err.Error())
		return
	}

	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(config)
	if err != nil {
		githubLogger.Panic("decode配置失败：" + err.Error())
		return
	}
}

func (p *github) loadReplaceStr(config githubConfig) {
	p.byteGithubReplace = fmt.Sprintf("%s://%s", config.ProxySiteScheme, config.ProxySiteDomain)
	p.byteApiReplace = fmt.Sprintf("%s://%s/~/api.github.com", config.ProxySiteScheme, config.ProxySiteDomain)
	p.byteRawReplace = fmt.Sprintf("%s://%s/~/raw.githubusercontent.com", config.ProxySiteScheme, config.ProxySiteDomain)
	p.stringObjectsContentReplace = fmt.Sprintf("%s://%s/~/objects.githubusercontent.com", config.ProxySiteScheme, config.ProxySiteDomain)

}

func (p *github) addToRegister() {
	proxySite := Register.NewProxySiteInfo()
	proxySite.Scheme = "https"
	proxySite.ResponseModify = p.ModifyResponse()
	proxySite.AutoGzip = true
	Register.AddProxySite("github.com", proxySite)
	Register.AddProxySite("api.github.com", proxySite)
	Register.AddProxySite("raw.githubusercontent.com", proxySite)
	Register.AddProxySite("objects.githubusercontent.com", proxySite)
	Register.AddMiddlewareFunc(p.RedirectGitClientMiddleware())

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
		//删除安全头防止无法加载
		delete(res.Header, "Content-Security-Policy")

		byteBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		byteBody = bytes.Replace(byteBody, []byte("https://github.com"), []byte(p.byteGithubReplace), -1)
		byteBody = bytes.Replace(byteBody, []byte("https://api.github.com"), []byte(p.byteApiReplace), -1)
		byteBody = bytes.Replace(byteBody, []byte("https://raw.githubusercontent.com"), []byte(p.byteRawReplace), -1)

		if res.Header.Get("Location") != "" {
			if strings.Index(res.Header.Get("Location"), "https://github.com") >= 0 {
				res.Header.Set("Location", strings.Replace(res.Header.Get("Location"), "https://github.com", string(p.byteGithubReplace), -1))
			} else if strings.Index(res.Header.Get("Location"), "https://objects.githubusercontent.com") >= 0 {
				res.Header.Set("Location", strings.Replace(res.Header.Get("Location"), "https://objects.githubusercontent.com", p.stringObjectsContentReplace, -1))
			} else if strings.Index(res.Header.Get("Location"), "https://raw.githubusercontent.com") >= 0 {
				res.Header.Set("Location", strings.Replace(res.Header.Get("Location"), "https://raw.githubusercontent.com", p.byteRawReplace, -1))
			} else {
				githubLogger.Errorf("出现未被记录的Location:%q,URL:%q", res.Header.Get("Location"), res.Request.RequestURI)
			}
		}

		//将修改后的body复制回body
		res.Body = io.NopCloser(bytes.NewReader(byteBody))
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
