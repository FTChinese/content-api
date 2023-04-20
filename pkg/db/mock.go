//go:build !production

package db

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/FTchinese/content-api/pkg/config"
	"github.com/jmoiron/sqlx"
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

func MockMySQL() ReadWriteMyDBs {
	MustSetupViper()
	return MustNewMyDBs()
}
