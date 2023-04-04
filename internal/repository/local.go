package repository

import "github.com/davfer/capi/pkg/model"

type RepositoryEntry struct {
	Name       string
	Collection model.Collection
}

type LocalRepository struct {
	items []RepositoryEntry
}

func NewLocalRepository(cols []model.Collection) *LocalRepository {
	items := make([]RepositoryEntry, len(cols))
	for i, col := range cols {
		items[i] = RepositoryEntry{
			Name:       col.Name,
			Collection: col,
		}
	}

	return &LocalRepository{
		items: items,
	}
}

func (r *LocalRepository) GetCollection(name string) (*model.Collection, error) {
	for _, item := range r.items {
		if item.Name == name {
			return &item.Collection, nil
		}
	}

	return nil, nil
}
