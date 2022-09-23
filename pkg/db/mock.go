//go:build !production

package db

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/content-api/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadConfigFile() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filepath.Join(home, "config", "env.dev.toml"))
}

func MustReadConfigFile() []byte {
	b, err := ReadConfigFile()
	if err != nil {
		panic(err)
	}

	return b
}

func MustSetupViper() {
	config.MustSetupViper(MustReadConfigFile())
}

func MockDB() *sqlx.DB {
	MustSetupViper()
	return MustNewMySQL(config.MustMySQLReadConn())
}
