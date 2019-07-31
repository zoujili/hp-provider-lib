package app

import (
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-core-go/pkg/v1/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "App provider test", test.LoadCustomReporters("../../test_provider_app.xml"))
}

var _ = Describe("App provider", func() {
	name := "TestApp"
	_ = os.Setenv("APP_NAME", name)
	// Proper testing for version is done in the core package.
	version.BuildString = "0.13.90 a722bdb 2018-01-09T22:32:37+01:00 go version go1.11 linux/amd64"

	It("Returns the app parameters", func() {
		c := NewConfigFromEnv()
		p := New(c)

		err := p.Init()
		Expect(err).ToNot(HaveOccurred())

		Expect(p.Name()).To(Equal(name))

		Expect(p.Version()).ToNot(BeNil())
		Expect(p.Version().String()).To(Equal(version.BuildString))

		Expect(p.Config.BasePath).To(Equal(defaultBasePath))
	})
})
