package datum

import "github.com/pkg/errors"

type DataKey struct {
	Typ   string
	Col   string
	Name  string
	Store string
	Ref   string
}

type DataContext struct {
	data        map[DataKey]string
	missingData []DataKey
}

func NewDataContext() *DataContext {
	return &DataContext{
		data: make(map[DataKey]string),
	}
}

func (c *DataContext) Set(key DataKey, value string) {
	for i, k := range c.missingData {
		if k.Name == key.Name && k.Col == key.Col {
			c.missingData = append(c.missingData[:i], c.missingData[i+1:]...)
			break
		}
	}

	c.data[key] = value
}

func (c *DataContext) Request(key DataKey) error {
	c.missingData = append(c.missingData, key)
	return nil
}

func (c *DataContext) Get(key DataKey) (string, error) {
	if val, ok := c.data[key]; ok {
		return val, nil
	}

	return "", errors.New("key not found")
}

func (c *DataContext) GetAll() map[DataKey]string {
	return c.data
}

func (c *DataContext) IsComplete() bool {
	return len(c.missingData) == 0
}

func (c *DataContext) GetMissingData() []DataKey {
	return c.missingData
}
