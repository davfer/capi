package finder

import (
	"github.com/davfer/capi/internal/repository"
	"github.com/davfer/capi/pkg/model"
)

type SimpleFinder struct {
	repo repository.Repository
}

func NewSimpleFinder(r repository.Repository) *SimpleFinder {
	return &SimpleFinder{
		repo: r,
	}
}

func (sr *SimpleFinder) IsANewPotentialCollection(input string) bool {
	return false
}

func (sr *SimpleFinder) FindPotentialCollection(input string) (*model.Collection, error) {
	// is a URL?
	//_, err := url.Parse(input)
	//if err != nil {
	//	return nil, err
	//}

	// is a collection?
	c, err := sr.repo.GetCollection(input)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (sr *SimpleFinder) FindPotentialResource(c *model.Collection, args []string) (*model.Resource, error) {
	// is a resource?
	for _, r := range c.Resources {
		if r.Name == args[1] {
			return &r, nil
		}
	}

	return nil, nil
}

func (sr *SimpleFinder) IsANewPotentialResource(c *model.Collection, args []string) bool {
	return false
}
