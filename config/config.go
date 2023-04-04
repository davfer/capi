package config

import "github.com/davfer/capi/pkg/model"

type CapiConfiguration struct {
	Collections []model.Collection `json:"collections" yaml:"collections"`
}
