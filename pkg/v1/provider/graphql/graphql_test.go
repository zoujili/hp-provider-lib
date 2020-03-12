package graphql

import (
	"encoding/json"
	"fmt"
	"github.azc.ext.hp.com/hp-business-platform/lib-core-go/pkg/v1/test"
	"github.azc.ext.hp.com/hp-business-platform/lib-provider-go/pkg/v1/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGraphQL(t *testing.T) {
	RegisterFailHandlerWithT(t, Fail)
	RunSpecsWithDefaultAndCustomReporters(t, "GraphQL provider test", test.LoadCustomReporters("../../test_provider_graphql.xml"))
}

var _ = Describe("GraphQL provider", func() {

	schema := `schema {
					query: Query
				}
				type Query{
					ping(): String!
				}`

	It("Starts the graphQL HTTP service", func() {
		logrus.SetLevel(logrus.DebugLevel)
		var p *GraphQL
		By("Creating and initializing the provider", func() {
			p = New(&Config{
				Port:             defaultPort,
				GraphiQLEnabled:  true,
				GraphiQLEndpoint: defaultGraphiQLEndpoint,
			})
			err := p.SetSchema(schema, &resolver{})
			Expect(err).NotTo(HaveOccurred())
			err = p.Init()
			Expect(err).NotTo(HaveOccurred())
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
		By("Performing a query", func() {
			query := `{"query": "{ping() {}}"}`
			resp, err := http.Post(fmt.Sprintf("http://localhost:%d", defaultPort), "application/json", strings.NewReader(query))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(200))

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			Expect(err).NotTo(HaveOccurred())

			Expect(data).To(HaveKey("data"))

			queryResponse := data["data"].(map[string]interface{})
			Expect(queryResponse).To(HaveKeyWithValue("ping", "pong"))
		})
	})
})

type resolver struct {
}

func (r *resolver) Ping() string {
	return "pong"
}
