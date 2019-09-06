package config

type BuildConfig struct {
	Production bool   `json:"production"`
	Version    string `json:"version"`
	BuildAt    string `json:"buildAt"`
}
