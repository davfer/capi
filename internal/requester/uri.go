package requester

import (
	"github.com/davfer/capi/internal/datum"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

var ErrUnsupportedUriType = errors.New("unsupported uri type")

type UriBuilder interface {
	Build(*model.Collection, *model.Resource) (string, error)
}

type BasicUriBuilder struct {
	dataSys *datum.System
}

func NewBasicUriBuilder(dataSys *datum.System) *BasicUriBuilder {
	return &BasicUriBuilder{
		dataSys: dataSys,
	}
}

func (b *BasicUriBuilder) Build(c *model.Collection, r *model.Resource) (string, error) {
	uri := &url.URL{}
	var err error
	if c.Base != "" {
		uri, err = url.Parse(c.Base)
		if err != nil {
			return "", errors.Wrap(err, "error parsing uri")
		}
	}

	uri.Path = uri.Path + "/" + strings.Trim(r.Path, "/")

	if _, ok := c.Environments["default"]; ok {
		if c.Environments["default"].Protocol != "" {
			uri.Scheme = c.Environments["default"].Protocol
		}

		if c.Environments["default"].Host != "" {
			uri.Host = c.Environments["default"].Host
		}
		if c.Environments["default"].Port != 0 && c.Environments["default"].Port != 443 && c.Environments["default"].Port != 80 {
			//port = strconv.Itoa(c.Environments["default"].Port)
		}
	}

	// query
	if c.Auth.Token != nil && c.Auth.Token.Place == "query" {
		name := "token"
		if c.Auth.Token.Var != "" {
			name = c.Auth.Token.Var
		}

		uri.Query().Add(name, c.Auth.Token.Token)
	}

	logrus.Debug("uri: ", uri.String())

	data, err := b.dataSys.Complete(uri.String(), c)
	if err != nil {
		return "", errors.Wrap(err, "error completing data in uri")
	}

	res, err := b.dataSys.Replace(uri.String(), data)
	if err != nil {
		return "", errors.Wrap(err, "error replacing data in uri")
	}

	return res, nil
}
