# go-kit学习指南 - 基础概念和架构


## 测试

uppercase
```bash
$ curl -i http://localhost:8080/uppercase?s=foo
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 18 Apr 2024 16:14:08 GMT
Content-Length: 12

{"v":"FOO"}

```

count
```bash
$ curl -i http://localhost:8080/count?s=foo
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 18 Apr 2024 16:11:47 GMT
Content-Length: 8

{"v":3}
```
