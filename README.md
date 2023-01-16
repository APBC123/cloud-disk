# CloudDisk
>轻量级云盘，基于go-zero,xorm实现

使用到的命令
```text
# 创建API服务
goctl api new core

# 启动服务
go run core.go -f etc/core-api.yaml

# 使用API文件生成代码
goctl api go -api core.api -dir . style go_zero
```
使用到的服务
```text
github.com/jordan-wright/email
#用于向用户邮箱发送验证码

go-redis
#用于邮箱注册时暂存验证码，用作比对

satori/go.uuid
#用于生成各表每条记录的identity

jwt-go
用于生成及解析token
```
腾讯云COS后台地址：https://console.cloud.tencent.com/cos/bucket

腾讯云COS帮助文档：https://cloud.tencent.com/document/product/436/31215
