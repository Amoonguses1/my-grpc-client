package port

import (
	"context"

	resl "github.com/amoonguses1/grpc-proto-study/protogen/go/resiliency"
	"google.golang.org/grpc"
)

type ResiliencyClientPort interface {
	UnaryResiliency(ctx context.Context, in *resl.ResiliencyRequest,
		opts ...grpc.CallOption) (*resl.ResiliencyResponse, error)
	ServerStreamingResiliency(ctx context.Context, in *resl.ResiliencyRequest,
		opts ...grpc.CallOption) (resl.ResiliencyService_ServerStreamingResiliencyClient, error)
	ClientStreamingResiliency(ctx context.Context,
		opts ...grpc.CallOption) (resl.ResiliencyService_ClientStreamingResiliencyClient, error)
	BiDirectionalResiliency(ctx context.Context,
		opts ...grpc.CallOption) (resl.ResiliencyService_BiDirectionalResiliencyClient, error)
}
