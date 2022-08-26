# TheresaProxyV2

基于go语言编写的可拓展反代服务器。默认内置反代github站点插件。

由于使用go语言构建，几乎无需运行环境，这也是和上一代基于nodejs构建最大的不同之一。

~~没人用就不好好写readme了~~

## 使用

1. 在releases页面直接下载对应版本程序
2. 创建`config/github.json`文件，填入下方内容

```json
{
  "proxy_site_scheme": "https",
  "proxy_site_domain": "example.com"
}
```

随后，在应用代理github站点时，会自动将部分github域名替换为https://example.com/DOMAIN`以允许用户正常访问GitHub。

目前可支持直接替换github域名来使用`git clone`通过http(s)协议克隆仓库。

可以经过本应用直接下载release页面下的非`source code`文件。

