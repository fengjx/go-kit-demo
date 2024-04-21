package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-kit/kit/auth/basic"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	svc := greetService{}

	// 中间件组合实现功能扩展
	helloEndpoint := makeHelloEndpoint(svc)
	helloEndpoint = basic.AuthMiddleware("foo", "bar", "hello")(helloEndpoint)
	helloEndpoint = timeMiddleware("hello")(helloEndpoint)

	satHelloHandler := httptransport.NewServer(
		helloEndpoint,
		decodeRequest,
		encodeResponse,
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			authorization := request.Header.Get("Authorization")
			// 增加一个拦截器
			ctx = context.WithValue(ctx, httptransport.ContextKeyRequestAuthorization, authorization)
			return ctx
		}),
	)

	http.Handle("/say-hello", satHelloHandler)
	log.Println("http server start")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type helloReq struct {
	Name string `json:"name"`
}

type helloResp struct {
	Msg string `json:"msg"`
}

type greetService struct {
}

func (svc greetService) SayHi(_ context.Context, name string) string {
	return fmt.Sprintf("hi: %s", name)
}

func makeHelloEndpoint(svc greetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*helloReq)
		msg := svc.SayHi(ctx, req.Name)
		return helloResp{
			Msg: msg,
		}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	req := &helloReq{
		Name: name,
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	data := map[string]any{
		"status": 0,
		"msg":    "ok",
		"data":   response,
	}
	return json.NewEncoder(w).Encode(data)
}

// timeMiddleware 打印耗时
func timeMiddleware(endpointName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			start := time.Now()
			defer func() {
				duration := time.Since(start)
				log.Printf("%s take %v\r\n", endpointName, duration)
			}()
			return next(ctx, request)
		}
	}
}
