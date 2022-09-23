package config

type ServerStatus struct {
	Version    string `json:"version"`
	BuiltAt    string `json:"builtAt"`
	Commit     string `json:"commit"`
	Production bool   `json:"production"`
}
