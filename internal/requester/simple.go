package requester

import (
	"bytes"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SimpleRequester struct {
	builder *HttpBuilderResolver
}

func NewSimpleRequester(builder *HttpBuilderResolver) *SimpleRequester {
	return &SimpleRequester{
		builder: builder,
	}
}

func (sr *SimpleRequester) New(c *model.Collection, r *model.Resource) *RequestCandidate {
	return &RequestCandidate{
		collection: c,
		resource:   r,
		builder:    sr.builder,
	}
}

type RequestCandidate struct {
	collection *model.Collection
	resource   *model.Resource
	builder    *HttpBuilderResolver

	request *http.Request
	client  http.Client
}

func (c *RequestCandidate) IsComplete() bool {
	return false
}

func (c *RequestCandidate) Execute() (*http.Response, error) {
	method := c.resource.Method
	uri, err := c.builder.BuildUri(c.collection, c.resource)
	if err != nil {
		return nil, errors.Wrap(err, "error building uri")
	}

	body, err := c.builder.BuildBody(c.collection, c.resource)
	if err != nil {
		return nil, errors.Wrap(err, "error building body")
	}

	headers, err := c.builder.BuildHeaders(c.collection, c.resource)
	if err != nil {
		return nil, errors.Wrap(err, "error building headers")
	}

	client, err := c.builder.BuildClient(c.collection, c.resource)
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}

	logrus.WithFields(logrus.Fields{
		"method":  method,
		"uri":     uri,
		"headers": headers,
	}).Info("executing request")

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return client.Do(req)
}
