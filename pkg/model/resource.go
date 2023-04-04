package model

type Resource struct {
	Name    string            `json:"name" yaml:"name"`
	Path    string            `json:"path" yaml:"path"`
	Method  string            `json:"method" yaml:"method"`
	Auth    string            `json:"auth" yaml:"auth"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Body    BodyTypes         `json:"body" yaml:"body"`
}

type BodyTypes struct {
	Json *BodyJson `json:"json" yaml:"json"`
}

type BodyJson map[string]any
