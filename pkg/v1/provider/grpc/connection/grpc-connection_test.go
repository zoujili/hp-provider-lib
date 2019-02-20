package connection

import (
	"context"
	"errors"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/examples/ping/server/gen"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestGRPCConnection(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "GRPC connection provider test", test.LoadCustomReporters("../../test_provider_grpc_connection.xml"))
}

var _ = Describe("GRPC connection provider test", func() {

	var server *grpc.Server

	BeforeSuite(func() {
		server = grpc.New(&grpc.Config{
			Port:       3030,
			LogPayload: true,
		})
		err := server.Init()
		Expect(err).NotTo(HaveOccurred())
		gen.RegisterPingServiceServer(server.Server, TestService{})
		go func() {
			err := server.Run()
			Expect(err).NotTo(HaveOccurred())
		}()
		err = provider.WaitForRunningProvider(server, 2*time.Second)
		Expect(err).NotTo(HaveOccurred())
		Expect(server.IsRunning()).To(BeTrue())
	})

	Context("The GRPC ping service is running", func() {
		It("Starts the GRPC connection", func() {
			logrus.SetLevel(logrus.DebugLevel)
			var p *Connection
			var client gen.PingServiceClient

			By("Creating and initializing the provider", func() {
				p = New(&Config{
					Host:       defaultHost,
					Port:       3030,
					LogPayload: true,
				})
				err := p.Init()
				Expect(err).NotTo(HaveOccurred())
			})

			By("Running the Connection", func() {
				go func() {
					err := p.Run()
					Expect(err).NotTo(HaveOccurred())
				}()
				err := provider.WaitForRunningProvider(p, 2*time.Second)
				Expect(err).NotTo(HaveOccurred())
				Expect(p.IsRunning()).To(BeTrue())
			})
			By("Creating the PingService client", func() {
				client = gen.NewPingServiceClient(p.Conn)
				Expect(client).NotTo(BeNil())
			})
			By("Calling the client", func() {
				res, err := client.Ping(context.Background(), &gen.PingRequest{In: "Hello"})
				Expect(err).NotTo(HaveOccurred())
				Expect(res).NotTo(BeNil())
				Expect(res.Out).To(Equal("Hello"))
			})
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
