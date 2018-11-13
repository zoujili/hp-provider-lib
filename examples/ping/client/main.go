package main

import (
	"context"
	"fitstation-hp/lib-fs-provider-go/examples/ping/server"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	stack := stack.New()
	defer stack.MustClose()

	logrusProvider := provider.NewLogrus(&provider.LogrusConfig{
		Level:     logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{},
		Output:    os.Stderr,
	})
	stack.MustInit(logrusProvider)

	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatal("did not connect")
	}
	defer conn.Close()

	client := server.NewPingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := "hello"
	if len(os.Args) > 1 {
		in = os.Args[1]
	}

	res, err := client.Ping(ctx, &server.PingRequest{In: in})
	if err != nil {
		logrus.WithError(err).Error("call failed")
	}

	fmt.Println(res.Out)
}
