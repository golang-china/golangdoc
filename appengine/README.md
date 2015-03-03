# GAE 版本 golangdoc

构建步骤:

1. 安装 [go_appengine](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go), 并添加到 `$PATH` 环境变量
2. 进入当前目录, 运行命令: `go run main.go`, 生成 `goroot.zip` 文件
3. 本地启动AGE程序: `goapp serve .`, 打开网页 http://127.0.0.1:8080
4. OK

部署到GAE:

1. 打开 `app.yaml` 文件, 设置 `application:` 字段为对应的 `app-id` (改成自己的 `APPID`)
2. 打开VPN, 运行 `goapp deploy` 上传应用
3. 打开网页 http://app-id.appspot.com/
4. OK

补充:
目前GAE版本没有启动搜索功能, 可以自己手工添加索引文件.
