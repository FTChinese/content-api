package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"testing"
)

var DB *sqlx.DB

func init() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var dbConn Conn
	err = viper.UnmarshalKey("mysql.dev", &dbConn)
	if err != nil {
		panic(err)
	}

	DB, err = NewDB(dbConn)
	if err != nil {
		panic(err)
	}
}

func TestEnv_RetrieveRawStory(t *testing.T) {

	env := Env{db: DB}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve raw story",
			args:    args{id: "001076308"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.RetrieveRawStory(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveRawStory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}
