# Hugo 阿里云OSS快捷同步

折腾Hugo博客的时候发现Hugo博客的上传没有Hexo那样方便.
找了下阿里云OSS的SDK API文档,发现有Go语言版本的,所以就整了个出来.

首次使用需要配置


基于 Golang1.19 开发

无需任何其他环境,打开即用,如果要修改配置文件,需要找到程序目录下的 config.json 修改 

## OSS SDK源码和API文档

请访问[这里](https://github.com/aliyun/aliyun-oss-go-sdk?spm=a2c4g.11186623.0.0.1e8017fbqQLyba)获取OSS Go SDK源码。更多信息请参见[OSS Go SDK API文档](https://pkg.go.dev/github.com/aliyun/aliyun-oss-go-sdk/oss)。 

OSS SDK 示例代码
[OSS Go SDK示例代码_对象存储 OSS-阿里云帮助中心](https://help.aliyun.com/document_detail/32144.html)

