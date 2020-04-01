>> 用于记录所有请求头以及请求body的插件，json格式。

filePath string 日志最终打印文件位置
bufferingSize string 缓存大小，不设置或者设置为0则不使用缓存, 支持K、M，最大16M
bodyEnable bool 是否记录body
bodyMaxSize  string body记录的最大byte，支持K、M，最大16M

使用示例：
spec:
  sfLogger:
    bodyEnable: true
    bodyMaxSize: 1K
    filePath: /tmp/test.log
    bufferingSize: 1K