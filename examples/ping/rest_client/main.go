package main

import (
	"fmt"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/provider"
	"github.azc.ext.hp.com/fitstation-hp/lib-fs-provider-go/pkg/v1/stack"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	logrusProvider := provider.NewLogrus(&provider.LogrusConfig{
		Level:     logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{},
		Output:    os.Stderr,
	})
	st.MustInit(logrusProvider)

	res, err := http.Get("http://127.0.0.1:8080/ping")
	if err != nil {
		logrus.WithError(err).Error("call failed")
	}

	fmt.Println(res.Status)
}
