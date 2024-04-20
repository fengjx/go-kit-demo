package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"

	"github.com/fengjx/go-kit-demo/greetgrpc/pb"
)

func main() {

	svc := greetService{}
	helloEndpoint := makeHelloEndpoint(svc)

	satHelloHandler := httptransport.NewServer(
		helloEndpoint,
		decodeRequest,
		encodeResponse,
	)

	startHTTPServer := func() {
		log.Println("http server start")
		http.Handle("/say-hello", satHelloHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
	go startHTTPServer()

	// grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
	pb.RegisterGreeterServer(grpcServer, newGreeterServer(helloEndpoint))
	startGRPCServer := func() {
		log.Println("grpc server start")
		ln, err := net.Listen("tcp", ":8000")
		if err != nil {
			panic(err)
		}
		log.Fatal(grpcServer.Serve(ln))
	}
	startGRPCServer()
}

type greetService struct {
}

func (svc greetService) SayHi(_ context.Context, name string) string {
	return fmt.Sprintf("hi: %s", name)
}

func makeHelloEndpoint(svc greetService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.HelloReq)
		msg := svc.SayHi(ctx, req.Name)
		return &pb.HelloResp{
			Msg: msg,
		}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	req := &pb.HelloReq{
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

type GreeterServer struct {
	pb.UnimplementedGreeterServer
	sayHello grpctransport.Handler
}

func newGreeterServer(e endpoint.Endpoint) pb.GreeterServer {
	svr := &GreeterServer{}
	svr.sayHello = grpctransport.NewServer(
		e,
		svr.decodeSayHello,
		svr.encodeSayHello,
	)
	return svr
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	_, resp, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HelloResp), nil
}

func (s *GreeterServer) decodeSayHello(_ context.Context, req interface{}) (interface{}, error) {
	helloReq := req.(*pb.HelloReq)
	return &pb.HelloReq{
		Name: helloReq.Name,
	}, nil
}

func (s *GreeterServer) encodeSayHello(_ context.Context, resp interface{}) (interface{}, error) {
	return resp, nil
}
