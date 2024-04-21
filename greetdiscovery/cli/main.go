package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		etcdServer = "192.168.1.200:2379"
		prefix     = "/services/greetsvc/"
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

	logger := kitlog.NewNopLogger()
	instancer, err := etcdv3.NewInstancer(cli, prefix, logger)
	if err != nil {
		panic(err)
	}

	factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {
		u, err := url.Parse(instance)
		if err != nil {
			return nil, nil, err
		}
		httpcli := httptransport.NewClient(http.MethodPost, u, encodeRequest, decodeResponse)
		return httpcli.Endpoint(), nil, nil
	}

	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(3, 3*time.Second, balancer)

	req := &helloReq{
		Name: "fengjx",
	}
	resp, err := retry(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resp: %v", resp)
}

type helloReq struct {
	Name string `json:"name"`
}

type helloResp struct {
	Msg string `json:"msg"`
}

type Data[T any] struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
}

// encodeRequest encodes the request into JSON format
func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(*helloReq)
	r.Method = http.MethodGet
	r.URL.Path = "/say-hello"
	q := r.URL.Query()
	q.Set("name", req.Name)
	r.URL.RawQuery = q.Encode()
	return nil
}

// decodeResponse decodes the response from JSON format
func decodeResponse(_ context.Context, response *http.Response) (interface{}, error) {
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("response code: %d\r\n", response.StatusCode)
	}
	data := &Data[helloResp]{}
	err := json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	return data.Data.Msg, nil
}
