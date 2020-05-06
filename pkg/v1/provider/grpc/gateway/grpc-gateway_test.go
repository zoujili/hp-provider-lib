package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/examples/ping/server/gen"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/app"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

const (
	defaultPortTest = 8888
)

func TestGRPCGateway(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "GRPC gateway provider test", test.LoadCustomReporters("../../test_provider_grpc_gateway.xml"))
}

var _ = Describe("GRPC gateway provider test", func() {

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

	It("Runs the GRPC gateway (no special settings)", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *Gateway

		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Port:       defaultPortTest,
				LogPayload: true,
				Enabled:    true,
			}, server, app.New(&app.Config{}))
			err := p.Init()
			Expect(err).NotTo(HaveOccurred())
		})
		By("Running the gateway", func() {
			go func() {
				defer GinkgoRecover()
				err := p.Run()
				Expect(err).NotTo(HaveOccurred())
			}()
			err := provider.WaitForRunningProvider(p, 2*time.Second)
			Expect(err).NotTo(HaveOccurred())
			Expect(p.IsRunning()).To(BeTrue())
		})
		By("Registering the gateway", func() {
			err := p.RegisterServices(gen.RegisterPingServiceHandler)
			Expect(err).NotTo(HaveOccurred())
		})
		By("Calling the gateway", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", defaultPortTest))
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.StatusCode).To(Equal(200))
		})
		By("Shutting down the gateway", func() {
			err := p.Close()
			Expect(err).ToNot(HaveOccurred())
			resetHTTPServer(p)
		})
	})
	It("Runs the GRPC gateway (with a basePath)", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *Gateway

		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Port:       defaultPortTest,
				LogPayload: true,
				Enabled:    true,
			}, server, app.New(&app.Config{
				BasePath: "/srv/api",
			}))
			err := p.Init()
			Expect(err).NotTo(HaveOccurred())
		})
		By("Running the gateway", func() {
			go func() {
				defer GinkgoRecover()
				err := p.Run()
				Expect(err).NotTo(HaveOccurred())
			}()
			err := provider.WaitForRunningProvider(p, 2*time.Second)
			Expect(err).NotTo(HaveOccurred())
			Expect(p.IsRunning()).To(BeTrue())
		})
		By("Registering the gateway", func() {
			err := p.RegisterServices(gen.RegisterPingServiceHandler)
			Expect(err).NotTo(HaveOccurred())
		})
		By("Calling the gateway", func() {
			res, err := http.Get(fmt.Sprintf("http://localhost:%d/srv/api/ping", defaultPortTest))
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.StatusCode).To(Equal(200))
		})
		By("Shutting down the gateway", func() {
			err := p.Close()
			Expect(err).ToNot(HaveOccurred())
			resetHTTPServer(p)
		})
	})
})

// Resets the shutdown state of the HTTP REST grpcSrv used by the Gateway, allowing it be be used again.
// It does this by using an unsafe pointer to the "inShutdown" field of http.Server and setting it to it's original value (0).
func resetHTTPServer(p *Gateway) {
	pv := reflect.ValueOf(p.srv)
	fv := reflect.Indirect(pv).FieldByName("inShutdown")
	fp := unsafe.Pointer(fv.UnsafeAddr())
	atomic.StoreInt32((*int32)(fp), 0)
}

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
