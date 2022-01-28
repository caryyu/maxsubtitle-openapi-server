![](https://travis-ci.org/caryyu/maxsubtitle-openapi-server.svg?branch=main) ![](https://img.shields.io/docker/pulls/caryyu/maxsubtitle-openapi-server.svg) 

# maxsubtitle-openapi-server

一个简易的字幕接口服务端，主要利用爬虫进行结构化数据处理并暴露接口出来给各方平台接入使用。

> 注意：所有的 `ass` 字幕会统一做 `srt` 转换处理，理由是 `ass` 中存在字体设置会存在各平台环境适配要求

## How to use

### Installation

```shell
go mod download
go run cmd/server/main.go
```

### Docker

```shell
docker run -d -p 3000:3000 caryyu/maxsubtitle-openapi-server:latest
```

## API

### 查询

GET http://localhost:3000/subtitle/search/mulan

```json
[
  {
    "id": "7474966f70dc5df50217aa73180f4a3a5f0aac1ea8d28a7bca740691e74d062c",
    "originalId": "122304",
    "desc": "花木兰 中文字幕 / 木兰传说 字幕下载 / 花木兰真人版 字幕 Mulan.2020.1080p.BluRay.x264.chs&eng.ass",
    "name": "default-122304.srt",
    "url": "https://www.a4k.net//system/files/subtitle/2020-12/a4k.net_1609230823.ass",
    "format": "srt"
  },
  {
    "id": "a82f6abd4c950a951de4d2f10560b5762a695acd5da168953cb842cafc1d77d9",
    "originalId": "119847",
    "desc": "花木兰 中文字幕 / Mulan 字幕下载 / 木兰传说 字幕 Mulan.2020.BluRay.chs.中影公映国配.srt",
    "name": "default-119847.srt",
    "url": "https://www.a4k.net//system/files/subtitle/2020-12/a4k.net_1607155518.srt",
    "format": "srt"
  }
]
```

### 下载

GET http://localhost:3000/subtitle/a82f6abd4c950a951de4d2f10560b5762a695acd5da168953cb842cafc1d77d9/download

```txt
...
Content-Disposition: attachment; filename=default-119847.srt
Content-Type: text/plain; charset=UTF-8
...
1
00:00:41,500 --> 00:00:45,150
有关花木兰的传说有很多
...
```

> 需要注意头信息 Content-Disposition ，该字段会告诉浏览器决定下载的文件名及后缀

# License

GPL 3.0
