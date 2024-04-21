# go-kit学习指南 - 中间件


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
