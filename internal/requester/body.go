package requester

import (
	"encoding/json"
	"github.com/davfer/capi/internal/datum"
	"github.com/davfer/capi/pkg/model"
	"github.com/pkg/errors"
)

var ErrUnsupportedBodyType = errors.New("unsupported body type")

type BodyBuilder interface {
	Build(*model.Collection, *model.Resource) ([]byte, error)
}

// JSON Body Builder
type JsonBodyBuilder struct {
	dataSys *datum.System
}

func NewJsonBodyBuilder(dataSys *datum.System) *JsonBodyBuilder {
	return &JsonBodyBuilder{
		dataSys: dataSys,
	}
}

func (j *JsonBodyBuilder) Build(c *model.Collection, r *model.Resource) ([]byte, error) {
	if r.Body.Json == nil {
		return nil, ErrUnsupportedBodyType
	}

	jsonString, err := json.Marshal(r.Body.Json)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling json body")
	}

	data, err := j.dataSys.Complete(string(jsonString), c)
	if err != nil {
		return nil, errors.Wrap(err, "error completing data in json body")
	}

	res, err := j.dataSys.Replace(string(jsonString), data)
	if err != nil {
		return nil, errors.Wrap(err, "error replacing data in json body")
	}

	return []byte(res), nil
}

// Nil Body Builder
type NilBodyBuilder struct {
}

func NewNilBodyBuilder() *NilBodyBuilder {
	return &NilBodyBuilder{}
}

func (n *NilBodyBuilder) Build(c *model.Collection, r *model.Resource) ([]byte, error) {
	return []byte{}, nil
}
