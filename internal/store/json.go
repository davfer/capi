package store

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

var ErrNotFoundKey = errors.New("not found")

type JsonStore struct {
	file string
	data map[string]map[string]string
}

func NewJsonStore(file string) *JsonStore {
	_, err := os.Stat(file)
	if err != nil {
		_, err = os.Create(file)
		if err != nil {
			return nil
		}
	}

	readFile, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	data := map[string]map[string]string{}
	err = json.Unmarshal(readFile, &data)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal json file: %s", err.Error()))
	}

	return &JsonStore{
		file: file,
		data: data,
	}
}

func (j *JsonStore) Get(col, key string) (string, error) {
	if j.data[col] == nil {
		return "", ErrNotFoundKey
	}

	return j.data[col][key], nil
}

func (j *JsonStore) Set(col, key, value string) error {
	if j.data[col] == nil {
		j.data[col] = map[string]string{}
	}

	j.data[col][key] = value

	data, err := json.Marshal(j.data)
	if err != nil {
		return err
	}

	err = os.WriteFile(j.file, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
