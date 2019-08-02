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

		Expect(p.ParseEndpoint()).To(Equal("/"))
		Expect(p.ParseEndpoint("sub", "elem")).To(Equal("/sub/elem"))
		Expect(p.ParsePath()).To(Equal("/"))
		Expect(p.ParsePath("sub", "elem")).To(Equal("/sub/elem/"))
	})

	Context("Configuring the base path", func() {
		It("Handles a path with suffixed and prefixed slash", func() {
			_ = os.Setenv("APP_BASE_PATH", "/some/path/")
			p := New(NewConfigFromEnv())

			err := p.Init()
			Expect(err).ToNot(HaveOccurred())

			Expect(p.Config.BasePath).To(Equal("/some/path"))
			Expect(p.ParseEndpoint()).To(Equal("/some/path"))
			Expect(p.ParseEndpoint("sub", "elem")).To(Equal("/some/path/sub/elem"))
			Expect(p.ParsePath()).To(Equal("/some/path/"))
			Expect(p.ParsePath("sub", "elem")).To(Equal("/some/path/sub/elem/"))
		})
		It("Handles a path without suffixed slash", func() {
			_ = os.Setenv("APP_BASE_PATH", "/some/path")
			p := New(NewConfigFromEnv())

			err := p.Init()
			Expect(err).ToNot(HaveOccurred())

			Expect(p.Config.BasePath).To(Equal("/some/path"))
			Expect(p.ParseEndpoint()).To(Equal("/some/path"))
			Expect(p.ParseEndpoint("sub", "elem")).To(Equal("/some/path/sub/elem"))
			Expect(p.ParsePath()).To(Equal("/some/path/"))
			Expect(p.ParsePath("sub", "elem")).To(Equal("/some/path/sub/elem/"))
		})
		It("Handles a path without prefixed slash", func() {
			_ = os.Setenv("APP_BASE_PATH", "some/path/")
			p := New(NewConfigFromEnv())

			err := p.Init()
			Expect(err).ToNot(HaveOccurred())

			Expect(p.Config.BasePath).To(Equal("/some/path"))
			Expect(p.ParseEndpoint()).To(Equal("/some/path"))
			Expect(p.ParseEndpoint("sub", "elem")).To(Equal("/some/path/sub/elem"))
			Expect(p.ParsePath()).To(Equal("/some/path/"))
			Expect(p.ParsePath("sub", "elem")).To(Equal("/some/path/sub/elem/"))
		})
	})
})
