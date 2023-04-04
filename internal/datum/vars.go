package datum

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type VarParser struct {
}

func NewVarParser() *VarParser {
	return &VarParser{}
}

func (p *VarParser) Parse(str string) ([]DataKey, error) {
	data := make([]DataKey, 0)

	regexRef := regexp.MustCompile(`\$_var\([a-zA-Z0-9@ ]+,[a-zA-Z0-9@ ]+,[a-zA-Z0-9@ ]+\)`)
	regexParts := regexp.MustCompile(`\$_var\(([a-zA-Z0-9@ ]+),([a-zA-Z0-9@ ]+),([a-zA-Z0-9@ ]+)\)`)

	matches := regexRef.FindAllString(str, -1)
	for _, match := range matches {
		parts := regexParts.FindAllStringSubmatch(match, -1)
		if len(parts) != 1 || len(parts[0]) != 4 {
			return nil, errors.New("invalid match")
		}

		collectionName := strings.Trim(parts[0][1], " ")
		resourceName := strings.Trim(parts[0][2], " ")
		storeName := strings.Trim(parts[0][3], " ")

		data = append(data, DataKey{
			Typ:   "var",
			Col:   collectionName,
			Name:  resourceName,
			Store: storeName,
			Ref:   match,
		})
	}

	return data, nil
}

func (p *VarParser) Replace(str string, key DataKey, value string) string {
	return strings.Replace(str, key.Ref, value, -1)
}
