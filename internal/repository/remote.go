package repository

import "github.com/davfer/capi/pkg/model"

type Repository interface {
	GetCollection(name string) (*model.Collection, error)
}

type RemoteRepository struct {
	url string
}

func NewRemoteRepository(url string) *RemoteRepository {
	return &RemoteRepository{
		url: url,
	}
}

func (r *RemoteRepository) GetCollection(name string) (*model.Collection, error) {
	return nil, nil
}
