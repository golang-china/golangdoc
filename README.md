# Godoc 改进版本, 支持翻译文档的动态加载

GAE预览 http://golang-china.appspot.com/

# 安装 golangdoc

安装 golangdoc :

	go get github.com/golang-china/golangdoc

下载翻译文件 到 `$(GOROOT)/translations` 目录:

	https://github.com/golang-china/golangdoc.translations

启用简体中文版文档服务:

	golangdoc -http=:6060 -lang=zh_CN

动态切换包文档:

	http://127.0.0.1:6060/pkg/builtin/
	http://127.0.0.1:6060/pkg/builtin/?lang=en
	http://127.0.0.1:6060/pkg/builtin/?lang=raw
	http://127.0.0.1:6060/pkg/builtin/?lang=zh_CN

其中 URL 的 `lang` 参数为 `en`/`raw` 或 无对应语言时 表示使用原始的文档,
缺少或为空时用 golangdoc 服务器启动时命令行指定的 `lang` 参数.

## 部署到 AGE 环境

golangdoc 支持 GAE 环境. 具体请参考: [appengine/README.md](appengine/README.md)


# 系统服务模式运行(Windows平台)

	# 安装 Windows 服务
	golangdoc -service-install -http=:6060

	# 启动/停止 Windows 服务
	golangdoc -service-start
	golangdoc -service-stop

	# 卸载 Windows 服务
	golangdoc -service-remove


# 其他

- GAE环境支持: https://github.com/golang-china/golangdoc/tree/master/appengine
- 文档翻译项目: http://github.com/golang-china/golangdoc.translations
- 文档提取工具: http://godoc.org/github.com/golang-china/golangdoc/docgen
- 本地化支持包: http://godoc.org/github.com/golang-china/golangdoc/local

