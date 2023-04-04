package requester

import (
	"github.com/davfer/capi/internal/datum"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
	"net/http"
)

type HttpBuilderResolver struct {
	BodyBuilders   []BodyBuilder
	UriBuilders    []UriBuilder
	HeaderBuilders []HeaderBuilder
	ClientBuilders []ClientBuilder
}

func NewHttpBuilderResolver(dataSys *datum.System) *HttpBuilderResolver {
	return &HttpBuilderResolver{
		BodyBuilders:   []BodyBuilder{NewJsonBodyBuilder(dataSys), NewNilBodyBuilder()},
		UriBuilders:    []UriBuilder{NewBasicUriBuilder(dataSys)},
		HeaderBuilders: []HeaderBuilder{NewBasicHttpHeaderBuilder(dataSys)},
		ClientBuilders: []ClientBuilder{NewDefaultClientBuilder(), NewSimpleClientBuilder()},
	}
}

func (h *HttpBuilderResolver) BuildBody(c *model.Collection, r *model.Resource) ([]byte, error) {
	for _, builder := range h.BodyBuilders {
		body, err := builder.Build(c, r)
		if errors.Is(err, ErrUnsupportedBodyType) {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "error building body")
		}

		return body, nil
	}

	return nil, errors.New("no body builder found")
}

func (h *HttpBuilderResolver) BuildUri(c *model.Collection, r *model.Resource) (string, error) {
	for _, builder := range h.UriBuilders {
		uri, err := builder.Build(c, r)
		if errors.Is(err, ErrUnsupportedUriType) {
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "error building uri")
		}

		return uri, nil
	}

	return "", errors.New("no uri builder found")
}

func (h *HttpBuilderResolver) BuildHeaders(c *model.Collection, r *model.Resource) (map[string]string, error) {
	headers := make(map[string]string)
	for _, builder := range h.HeaderBuilders {
		hdrs, err := builder.Build(c, r)
		if errors.Is(err, ErrUnsupportedHeaderType) {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "error building headers")
		}

		for k, v := range hdrs {
			headers[k] = v
		}
	}

	return headers, nil
}

func (h *HttpBuilderResolver) BuildClient(c *model.Collection, r *model.Resource) (*http.Client, error) {
	for _, builder := range h.ClientBuilders {
		client, err := builder.Build(c, r)
		if errors.Is(err, ErrUnsupportedClientType) {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "error building client")
		}

		return client, nil
	}

	return nil, errors.New("no client builder found")
}
