package main

import (
	"context"
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/examples/ping/server"
	provider "github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/logrus"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	logrusProvider := provider.New(&provider.Config{
		Level:     logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{},
		Output:    os.Stderr,
	})
	st.MustInit(logrusProvider)

	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		logrus.WithError(err).Fatal("did not connect")
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.WithError(err).Error("error while closing connection")
		}
	}()

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
