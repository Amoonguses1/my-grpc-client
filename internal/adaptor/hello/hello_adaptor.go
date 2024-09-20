package hello

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
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

func (a *HelloAdaptor) SayManyHellos(ctx context.Context, name string) {
	helloRequest := &hello.HelloRequest{
		Name: name,
	}

	greetStream, err := a.helloClient.SayManyHellos(ctx, helloRequest)
	if err != nil {
		log.Fatalln("Error on SayHello: ", err)
	}

	for {
		greet, err := greetStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error on SayManyHellos: ", err)
		}

		log.Println(greet.Greet)
	}
}

func (a *HelloAdaptor) SayHelloToEveryone(ctx context.Context, names []string) {
	greetStream, err := a.helloClient.SayHelloToEveryone(ctx)
	if err != nil {
		log.Fatalln("Error on SayHelloToEveryone :", err)
	}

	for _, name := range names {
		req := &hello.HelloRequest{
			Name: name,
		}

		greetStream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := greetStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error on SayHelloToEveryone :", err)
	}

	log.Println(res.Greet)
}

func (a *HelloAdaptor) SayHelloContinuous(ctx context.Context, names []string) {
	greetStream, err := a.helloClient.SayHelloContinuous(ctx)
	if err != nil {
		log.Fatalln("Error on SayHelloContinuous :", err)
	}

	greetChan := make(chan struct{})

	go func() {
		for _, name := range names {
			req := &hello.HelloRequest{
				Name: name,
			}

			greetStream.Send(req)
		}

		greetStream.CloseSend()
	}()

	go func() {
		for {
			greet, err := greetStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("Error on SayHelloContinuous :", err)
			}

			log.Println(greet.Greet)
		}

		close(greetChan)
	}()

	<-greetChan
}
