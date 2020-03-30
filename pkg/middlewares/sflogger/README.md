>> 用于记录所有请求头以及请求body的插件，json格式。

filePath string 日志最终打印文件位置
bufferingSize int64 缓存大小，不设置或者设置为0则不使用缓存 
bodyEnable bool 是否记录body
bodyMaxSize  int64 body记录的最大byte

使用示例：
spec:
  sfLogger:
    bodyEnable: true
    bodyMaxSize: 1048576
    filePath: /tmp/test.log
    bufferingSize: 1024