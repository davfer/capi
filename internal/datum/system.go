package datum

import (
	"errors"
	"github.com/davfer/capi/internal/store"
	"github.com/davfer/capi/pkg/model"
)

var ErrIncompleteData = errors.New("incomplete data")

type System struct {
	parsers  map[string]Parser
	resolver Resolver
}

func NewSystem() *System {
	return &System{
		parsers:  map[string]Parser{"var": NewVarParser()},
		resolver: NewDbResolver(store.NewJsonStore("./data.json"), NewAskResolver()),
	}
}

func (s *System) Complete(str string, col *model.Collection) (*DataContext, error) {
	data := NewDataContext()

	for _, parser := range s.parsers {
		keys, err := parser.Parse(str)
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			if key.Col == "@self" {
				key.Col = col.Name
			}

			err = data.Request(key)
			if err != nil {
				return nil, err
			}
		}
	}

	err := s.resolver.Resolve(data, col)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *System) Replace(str string, data *DataContext) (string, error) {
	if !data.IsComplete() {
		return "", ErrIncompleteData
	}

	res := str
	for k, v := range data.GetAll() {
		parser, ok := s.parsers[k.Typ]
		if !ok {
			continue
		}

		res = parser.Replace(res, k, v)
	}

	return res, nil
}
