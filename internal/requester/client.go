package requester

import (
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var ErrUnsupportedClientType = errors.New("unsupported client type")

type ClientBuilder interface {
	Build(*model.Collection, *model.Resource) (*http.Client, error)
}

type SimpleClientBuilder struct {
}

func NewSimpleClientBuilder() *SimpleClientBuilder {
	return &SimpleClientBuilder{}
}

func (s *SimpleClientBuilder) Build(c *model.Collection, r *model.Resource) (*http.Client, error) {
	return &http.Client{
		Timeout: 10 * time.Second,
	}, nil
}

type DefaultClientBuilder struct {
}

func NewDefaultClientBuilder() *DefaultClientBuilder {
	return &DefaultClientBuilder{}
}

func (d *DefaultClientBuilder) Build(c *model.Collection, r *model.Resource) (*http.Client, error) {
	return http.DefaultClient, nil
}
