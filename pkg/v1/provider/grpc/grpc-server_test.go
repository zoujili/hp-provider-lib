package grpc

import (
	"context"
	"errors"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/examples/ping/server/gen"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestGRPCServer(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "GRPC server provider test", test.LoadCustomReporters("../../test_provider_grpc_server.xml"))
}

var _ = Describe("GRPC server provider test", func() {
	It("Starts the GRPC server", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *Server

		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Port:         defaultPort,
				LogPayload:   true,
				EnableHealth: true,
			})
			err := p.Init()
			Expect(err).NotTo(HaveOccurred())
		})
		By("Registering the TestService", func() {
			gen.RegisterPingServiceServer(p.Server, TestService{})
		})
		By("Running the provider", func() {
			go func() {
				err := p.Run()
				Expect(err).NotTo(HaveOccurred())
			}()
			err := provider.WaitForRunningProvider(p, 2*time.Second)
			Expect(err).NotTo(HaveOccurred())
			Expect(p.IsRunning()).To(BeTrue())
		})
		By("Dialing into the grpc service", func() {
			conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
			Expect(err).NotTo(HaveOccurred())
			Expect(conn).NotTo(BeNil())

			request := gen.PingRequest{
				In: "Hello",
			}

			response := gen.PingResponse{}

			err = conn.Invoke(context.Background(), "/api.PingService/Ping", &request, &response)
			Expect(err).NotTo(HaveOccurred())
		})
		By("Shutting down the server", func() {
			err := p.Close()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

type TestService struct {
}

func (s TestService) Ping(ctx context.Context, request *gen.PingRequest) (*gen.PingResponse, error) {
	logger := ctxlogrus.Extract(ctx)
	logger.Info("hello from ping")

	if request.In == "panic" {
		panic("please panic")
	}

	if request.In == "error" {
		return nil, errors.New("please error me")
	}

	return &gen.PingResponse{Out: request.In}, nil
}
