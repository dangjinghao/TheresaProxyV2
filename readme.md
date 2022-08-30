# TheresaProxyV2

基于go语言编写的可拓展反代服务器。默认内置反代github反代插件。

由于使用go语言构建，几乎无需任何环境依赖，这也是和上一代基于nodejs构建最大的不同之一。

建议仅作为浏览使用，不要登陆

~~star能到两位数就开个项目页面~~

## DEMO

https://proxy.qacgn.com/github.com

## Bench mark

没有基准测试，性能预估十分捉急。

## 安装

1. 从本仓库的release页面下载对应操作系统，版本的程序

2. 在程序所在目录下创建`config/github.json`，并填入以下内容

   ```json
   {
     "proxy_site_scheme": "SCHEME",
     "proxy_site_domain": "YOUR_DOMAIN"
   }
   ```

3. 为下载的程序添加执行权限，运行程序

## 从源码构建

clone本代码，直接构建

## 使用

### GITHUB反代

**特性支持列表**

- [x] 几乎无感的正常浏览（~~能用就行~~）
- [x] 反代`objects.githubusercontent.com`实现部分文件反代下载
- [x] git clone http(s)支持
- [ ] **git clone ssh 不支持**（无计划）

#### 浏览github站点

直接访问`SCHEME://YOUR_DOMAIN/github.com`

#### 下载文件

- 当通过浏览器浏览时，可直接点击对应链接下载。
- 当使用curl/wget等工具时，将`https://github.com/xxx/xxx`替换为`SCHEME://YOUR_DOMAIN/~/github.com/xxx/xxx`。例如`wget https://github.com/dangjinghao/TheresaProxyV2/raw/master/readme.md`替换为`wget https://proxy.qacgn.com/~/github.com/dangjinghao/TheresaProxyV2/raw/master/readme.md`

#### git http(s) 克隆

当使用`git clone https://github.com/xxx/xxx `时，直接修改为`git clone SCHEME://YOUR_DOMAIN/xxx/xxx`即可。例如`git clone https://github.com/dangjinghao/TheresaProxyV2` 可直接修改为`git clone https://proxy.qacgn.com/dangjinghao/TheresaProxyV2`。

## 插件开发

*TODO*
