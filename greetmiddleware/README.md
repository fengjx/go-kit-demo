# go-kit学习指南 - 中间件

文档：<https://blog.fengjx.com/pages/d6f092/>

## 启动服务

```bash
go run main.go
```

## 测试

认证失败
```bash
# sya-hello
curl http://localhost:8080/say-hello?name=fengjx
```

认证成功
```bash
curl -H 'Authorization: Basic Zm9vOmJhcg==' http://localhost:8080/say-hello?name=fengjx
```


## 相关项目

- [luchen](https://github.com/fengjx/luchen) 基于go-kit封装的微服务框架
- [lca](https://github.com/fengjx/lca) 基于 amis 实现的低代码后台系统
- [glca](https://github.com/fengjx/glca) lca 接口实现。基于`luchen`框架开发
- [lc](https://github.com/fengjx/lc) glca 的命令行工具

