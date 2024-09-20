package port

import (
	"context"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClient interface {
	SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error)
}
