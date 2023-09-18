# go-mitm

中间人代理。可以更方便地查看请求/响应结果，可以将请求转化成常用语言的请求代码。

[go-mitm](https://github.com/lizongying/go-mitm)

## Usage

1. 当前设备安装并信任CA证书 [ca.crt](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/ca.crt)
   ![image](./screenshot/img_5.png)

2. 执行程序
   mac: [mitm_darwin_arm64](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/mitm_darwin_arm64) [mitm_darwin_amd64](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/mitm_darwin_amd64)

   linux: [mitm_linux_arm64](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/mitm_linux_arm64) [mitm_linux_amd64](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/mitm_linux_amd64)

   windows: [mitm_windows_amd64.exe](https://github.com/lizongying/go-mitm/releases/download/v0.0.2/mitm_windows_amd64.exe)

3. 设置代理为 http://localhost:8082
   ![image](./screenshot/img_6.png)

4. 访问界面 http://localhost:8083

   ![image](./screenshot/img.png)
   ![image](./screenshot/img_1.png)
   ![image](./screenshot/img_4.png)
   ![image](./screenshot/img_2.png)
   ![image](./screenshot/img_3.png)

## Dev

### Test

```shell
go run ./cmd/mitm-web/*.go

# http://localhost:5173/
npm run --prefix ./web/ui dev

curl -X POST "https://httpbin.org/post" -H "accept: application/json" --data '{"a":"xyz","b":"123"}' -x http://localhost:8082 --cacert ./static/tls/ca.crt

```

### Build

```shell
make
```

## TODO

* python http.client
* node https/request
* java
* middleware tmpdir save image/video
