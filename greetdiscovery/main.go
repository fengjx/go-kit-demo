package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		etcdServer = "192.168.1.200:2379"
		prefix     = "/services/greetsvc/"
		instance   = "192.168.1.163:8080"
		key        = prefix + instance
		value      = "http://" + instance
		ctx        = context.Background()
	)
	// 创建 etcd 客户端连接
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second,
		DialKeepAlive: time.Second * 30,
	}
	cli, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		log.Panic(err)
	}
	// 创建日志记录器
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	// 创建注册与发现中间件
	r := etcdv3.NewRegistrar(cli, etcdv3.Service{
		Key:   key,
		Value: value,
	}, logger)
	r.Register()
	defer r.Deregister()

	svc := greetService{}
	satHelloHandler := httptransport.NewServer(
		makeHelloEndpoint(svc),
		decodeRequest,
		encodeResponse,
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
