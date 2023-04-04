package requester

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davfer/capi/internal/datum"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
)

var ErrUnsupportedHeaderType = errors.New("unsupported header type")

type HeaderBuilder interface {
	Build(*model.Collection, *model.Resource) (map[string]string, error)
}

type BasicHttpHeaderBuilder struct {
	dataSys *datum.System
}

func NewBasicHttpHeaderBuilder(dataSys *datum.System) *BasicHttpHeaderBuilder {
	return &BasicHttpHeaderBuilder{
		dataSys: dataSys,
	}
}

func (b *BasicHttpHeaderBuilder) Build(c *model.Collection, r *model.Resource) (map[string]string, error) {
	// declared headers
	var baseHeaders, envHeaders, resourceHeaders map[string]string

	if c.Headers != nil {
		baseHeaders = c.Headers
	}
	if _, ok := c.Environments["default"]; ok {
		if c.Environments["default"].Headers != nil {
			envHeaders = c.Environments["default"].Headers
		}
	}
	if r.Headers != nil {
		resourceHeaders = r.Headers
	}

	headers := make(map[string]string)
	for k, v := range baseHeaders {
		headers[k] = v
	}
	for k, v := range envHeaders {
		headers[k] = v
	}
	for k, v := range resourceHeaders {
		headers[k] = v
	}

	// auth headers
	if c.Auth.Basic != nil {
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Auth.Basic.Username, c.Auth.Basic.Password)))
	}
	if c.Auth.Token != nil && c.Auth.Token.Place == "header" {
		name := "Authorization"
		if c.Auth.Token.Var != "" {
			name = c.Auth.Token.Var
		}
		headers[name] = c.Auth.Token.Token
	}

	// body type headers
	if r.Body.Json != nil {
		headers["Content-Type"] = "application/json"
	}

	// parse and complete
	jsonString, err := json.Marshal(headers)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling headers")
	}

	data, err := b.dataSys.Complete(string(jsonString), c)
	if err != nil {
		return nil, errors.Wrap(err, "error completing headers")
	}

	res, err := b.dataSys.Replace(string(jsonString), data)
	if err != nil {
		return nil, errors.Wrap(err, "error replacing data in headers")
	}

	mapBack := make(map[string]string)
	err = json.Unmarshal([]byte(res), &mapBack)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling back headers")
	}

	return mapBack, nil
}
