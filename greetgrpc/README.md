# go-kit学习指南 - 多协议支持

文档：<https://blog.fengjx.com/pages/efb886/>

## 安装编译工具

- protoc 安装：<https://grpc.io/docs/protoc-installation/>
- protoc grpc 插件安装：<https://grpc.io/docs/languages/go/quickstart/> 


## 编译 proto 文件

```bash
cd pb
bash build.sh
```

编译后会生成 `greet.pb.go` 和 `greet_grpc.pb.go` 两个文件
```bash
$ ls    
build.sh         greet.pb.go      greet.proto      greet_grpc.pb.go
```

## 启动服务

```bash
go run main.go
```

## 测试


http 协议
```bash
# sya-hello
curl http://localhost:8080/say-hello?name=fengjx

{"data":{"msg":"hi: fengjx"},"msg":"ok","status":0}
```

grpc 协议
```bash
go run cmd/greetcli/mian.go

2024/04/20 15:46:10 hi: fengjx 
```


## 相关项目

- [luchen](https://github.com/fengjx/luchen) 基于go-kit封装的微服务框架
- [lca](https://github.com/fengjx/lca) 基于 amis 实现的低代码后台系统
- [glca](https://github.com/fengjx/glca) lca 接口实现。基于`luchen`框架开发
- [lc](https://github.com/fengjx/lc) glca 的命令行工具


