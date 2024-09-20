package resiliency

import (
	"context"
	"io"
	"log"

	resl "github.com/amoonguses1/grpc-proto-study/protogen/go/resiliency"
	"github.com/amoonguses1/my-grpc-client/internal/port"
	"google.golang.org/grpc"
)

type ResiliencyAdaptor struct {
	resiliencyClient             port.ResiliencyClientPort
	resiliencyWithMetadataClient port.ResiliencyWithMetadataClientPort
}

func NewResiliencyAdaptor(conn *grpc.ClientConn) (*ResiliencyAdaptor, error) {
	client := resl.NewResiliencyServiceClient(conn)
	clientWithMetadata := resl.NewResiliencyWithMetadataServiceClient(conn)

	return &ResiliencyAdaptor{
		resiliencyClient:             client,
		resiliencyWithMetadataClient: clientWithMetadata,
	}, nil
}

func (a *ResiliencyAdaptor) UnaryResiliency(ctx context.Context, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) (*resl.ResiliencyResponse, error) {
	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCodes:    statusCodes,
	}

	res, err := a.resiliencyClient.UnaryResiliency(ctx, resiliencyRequest)
	if err != nil {
		log.Println("Error on UnaryResiliency :", err)
		return nil, err
	}

	return res, nil
}

func (a *ResiliencyAdaptor) ServerStreamingResiliency(ctx context.Context, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32) {
	resiliencyRequest := &resl.ResiliencyRequest{
		MinDelaySecond: minDelaySecond,
		MaxDelaySecond: maxDelaySecond,
		StatusCodes:    statusCodes,
	}

	reslStream, err := a.resiliencyClient.ServerStreamingResiliency(ctx, resiliencyRequest)

	if err != nil {
		log.Fatalln("Error on ServerStreamingResiliency :", err)
	}

	for {
		res, err := reslStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error on ServerStreamingResiliency :", err)
		}

		log.Println(res.DummyString)
	}
}

func (a *ResiliencyAdaptor) ClientStreamingResiliency(ctx context.Context, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32, count int) {
	reslStream, err := a.resiliencyClient.ClientStreamingResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliency :", err)
	}

	for i := 1; i <= count; i++ {
		resiliencyRequest := &resl.ResiliencyRequest{
			MinDelaySecond: minDelaySecond,
			MaxDelaySecond: maxDelaySecond,
			StatusCodes:    statusCodes,
		}

		reslStream.Send(resiliencyRequest)
	}

	res, err := reslStream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error on ClientStreamingResiliency :", err)
	}

	log.Println(res.DummyString)
}

func (a *ResiliencyAdaptor) BiDirectionalResiliency(ctx context.Context, minDelaySecond int32,
	maxDelaySecond int32, statusCodes []uint32, count int) {
	reslStream, err := a.resiliencyClient.BiDirectionalResiliency(ctx)
	if err != nil {
		log.Fatalln("Error on BiDirectionalResiliency :", err)
	}

	reslChan := make(chan struct{})

	go func() {
		for i := 1; i <= count; i++ {
			resiliencyRequest := &resl.ResiliencyRequest{
				MinDelaySecond: minDelaySecond,
				MaxDelaySecond: maxDelaySecond,
				StatusCodes:    statusCodes,
			}

			reslStream.Send(resiliencyRequest)
		}

		reslStream.CloseSend()
	}()

	go func() {
		for {
			res, err := reslStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("Error on BiDirectionalResiliency :", err)
			}

			log.Println(res.DummyString)
		}

		close(reslChan)
	}()

	<-reslChan
}
