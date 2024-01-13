package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/go-kit-demo/greetsvc/endpoint"
	"github.com/fengjx/go-kit-demo/greetsvc/service"
	"github.com/fengjx/go-kit-demo/greetsvc/transport"
)

func main() {
	svc := service.NewAddSvc()
	sumHandler := httptransport.NewServer(
		endpoint.MakeSumEndpoint(svc),
		transport.DecodeSumRequest,
		transport.EncodeResponse,
	)

	concatHandler := httptransport.NewServer(
		endpoint.MakeConcatEndpoint(svc),
		transport.DecodeCountRequest,
		transport.EncodeResponse,
	)

	http.Handle("/sum", sumHandler)
	http.Handle("/concat", concatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
