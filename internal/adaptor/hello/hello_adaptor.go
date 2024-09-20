package hello

import (
	"context"
	"log"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
	"github.com/amoonguses1/my-grpc-client/internal/adaptor/hello"
	"github.com/amoonguses1/my-grpc-client/internal/port"
	"google.golang.org/grpc"
)

type HelloAdaptor struct {
	helloClient port.HelloClientPort
}

func NewHelloAdaptor(conn *grpc.ClientConn) (*HelloAdaptor, error) {
	client := hello.NewHelloServiceClient(conn)

	return &HelloAdaptor{
		helloClient: client,
	}, nil
}

func (a *HelloAdaptor) SayHello(ctx context.Context, name string) (*hello.HelloResponse, error) {
	helloRequest := &hello.HelloRequest{
		Name: name,
	}

	greet, err := a.helloClient.SayHello(ctx, helloRequest)
	if err != nil {
		log.Fatalln("Error on SayHello: ", err)
	}

	return greet, nil
}
