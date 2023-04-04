package model

type Collection struct {
	Name         string                 `json:"name" yaml:"name"`
	Environments map[string]Environment `json:"environments" yaml:"environments"`
	Resources    []Resource             `json:"resources" yaml:"resources"`
	Base         string                 `json:"base" yaml:"base"`
	Auth         Auth                   `json:"auth" yaml:"auth"`
	Headers      map[string]string      `json:"headers" yaml:"headers"`
}

type Environment struct {
	Protocol  string            `json:"protocol" yaml:"protocol"`
	Host      string            `json:"host" yaml:"host"`
	Port      int               `json:"port" yaml:"port"`
	Confirm   string            `json:"confirm" yaml:"confirm"`
	Message   string            `json:"message" yaml:"message"`
	Cookies   bool              `json:"cookies" yaml:"cookies"`
	Auth      Auth              `json:"auth" yaml:"auth"`
	Headers   map[string]string `json:"headers" yaml:"headers"`
	Secrets   Secrets           `json:"secrets" yaml:"secrets"`
	Variables map[string]string `json:"variables" yaml:"variables"`
}

type Secrets struct {
	Sops SopsSecret `json:"sops" yaml:"sops"`
}

type SopsSecret struct {
	File string `json:"file" yaml:"file"`
}

type Auth struct {
	Basic *BasicAuth                `json:"basic" yaml:"basic"`
	Token *TokenAuth                `json:"token" yaml:"token"`
	From  struct{ Resource string } `json:"from" yaml:"from"`
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type TokenAuth struct {
	Place string `json:"place" yaml:"place"`
	Var   string `json:"var" yaml:"var"`
	Token string `json:"token" yaml:"token"`
}
