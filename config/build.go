package config

type BuildConfig struct {
	Production bool   `json:"production"`
	Version    string `json:"version"`
	BuiltAt    string `json:"builtAt"`
	Commit     string `json:"commit"`
}
