package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/amoonguses1/my-grpc-client/internal/adaptor/bank"
	"github.com/amoonguses1/my-grpc-client/internal/adaptor/hello"
	dbank "github.com/amoonguses1/my-grpc-client/internal/application/domain/bank"
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

	// helloAdaptor, err := hello.NewHelloAdaptor(conn)
	// if err != nil {
	// 	log.Fatalln("Can not create HelloAdaptor :", err)
	// }

	bankAdaptor, err := bank.NewBankAdapter(conn)
	if err != nil {
		log.Fatalln("Can not create BankAdapter :", err)
	}

	// runSayHello(helloAdaptor, "my name")
	// runSayManyHellos(helloAdaptor, "call multiple name")
	// runSayHelloToEveryone(helloAdaptor, []string{"Andy", "Bob", "Chris"})
	// runSayHelloContinuous(helloAdaptor, []string{"Andy", "Bob", "Chris"})
	// runGetCurrentBalance(bankAdaptor, "7835697001xxxxx")
	// runFetchExchangeRates(bankAdaptor, "USD", "JPN")
	// runSummarizeTransactions(bankAdaptor, "7835697002yyyyy", 10)
	runTransferMultiple(bankAdaptor, "7835697004", "7835697003", 200)
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

func runGetCurrentBalance(adaptor *bank.BankAdapter, acct string) {
	bal, err := adaptor.GetCurrentBalance(context.Background(), acct)

	if err != nil {
		log.Fatalln("Failed to call GetCurrentBalance :", err)
	}

	log.Println(bal)
}

func runFetchExchangeRates(adaptor *bank.BankAdapter, fromCur string, toCur string) {
	adaptor.FetchExchangeRates(context.Background(), fromCur, toCur)
}

func runSummarizeTransactions(adaptor *bank.BankAdapter, acct string, numDummyTransactions int) {
	var tx []dbank.Transaction

	for i := 1; i <= numDummyTransactions; i++ {
		ttype := dbank.TransactionTypeIn

		if i%3 == 0 {
			ttype = dbank.TransactionTypeOut
		}

		t := dbank.Transaction{
			Amount:          float64(rand.Intn(500) + 10),
			TransactionType: ttype,
			Notes:           fmt.Sprintf("Dummy transaction %v", i),
		}

		tx = append(tx, t)
	}

	adaptor.SummarizeTransactions(context.Background(), acct, tx)
}

func runTransferMultiple(adaptor *bank.BankAdapter, fromAcct string, toAcct string,
	numDummyTransactions int) {
	var trf []dbank.TransferTransaction

	for i := 1; i <= numDummyTransactions; i++ {
		tr := dbank.TransferTransaction{
			FromAccountNumber: fromAcct,
			ToAccountNumber:   toAcct,
			Currency:          "USD",
			Amount:            float64(rand.Intn(200) + 5),
		}

		trf = append(trf, tr)
	}

	adaptor.TransferMultiple(context.Background(), trf)
}
