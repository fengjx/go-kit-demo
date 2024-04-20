package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fengjx/go-kit-demo/greetgrpc/pb"
)

func main() {
	clientConn, err := grpc.Dial(
		"localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	greeterClient := pb.NewGreeterClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	helloResp, err := greeterClient.SayHello(ctx, &pb.HelloReq{
		Name: "fengjx",
	})
	fmt.Printf("message -> %s\r\n", helloResp.Msg)
}
