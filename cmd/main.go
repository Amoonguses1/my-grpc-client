package main

import (
	"context"
	"log"

	"github.com/amoonguses1/my-grpc-client/internal/adaptor/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		log.Fatalln("Can not connect to gRPC server :", err)
	}
	defer conn.Close()

	helloAdaptor, err := hello.NewHelloAdaptor(conn)
	if err != nil {
		log.Fatalln("Can not create HelloAdaptor :", err)
	}

	// runSayHello(helloAdaptor, "my name")
	// runSayManyHellos(helloAdaptor, "call multiple name")
	// runSayHelloToEveryone(helloAdaptor, []string{"Andy", "Bob", "Chris"})
	runSayHelloContinuous(helloAdaptor, []string{"Andy", "Bob", "Chris"})
}

func runSayHello(adaptor *hello.HelloAdaptor, name string) {
	greet, err := adaptor.SayHello(context.Background(), name)
	if err != nil {
		log.Fatalln("Can not call sayHello")
	}

	log.Println(greet.Greet)
}

func runSayManyHellos(adaptor *hello.HelloAdaptor, name string) {
	adaptor.SayManyHellos(context.Background(), name)
}

func runSayHelloToEveryone(adaptor *hello.HelloAdaptor, names []string) {
	adaptor.SayHelloToEveryone(context.Background(), names)
}

func runSayHelloContinuous(adaptor *hello.HelloAdaptor, names []string) {
	adaptor.SayHelloContinuous(context.Background(), names)
}
