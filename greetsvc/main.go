package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	hello2 "github.com/fengjx/go-kit-demo/greetsvc/service/hello"
	"github.com/fengjx/go-kit-demo/greetsvc/transport"
)

func main() {
	svc := hello2.NewAddSvc()
	sumHandler := httptransport.NewServer(
		hello2.MakeSumEndpoint(svc),
		transport.DecodeSumRequest,
		transport.EncodeResponse,
	)

	concatHandler := httptransport.NewServer(
		hello2.MakeConcatEndpoint(svc),
		transport.DecodeCountRequest,
		transport.EncodeResponse,
	)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	log.Println("http server start")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
