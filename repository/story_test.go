package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/content-api/models"
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

	env := ContentEnv{db: DB}

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

			t.Logf("Raw story: %+v\n", got)

			t.Logf("CN story: %+v", models.NewStoryCN(&got))

			storyEN, err := models.NewStoryEN(&got)
			if err != nil {
				t.Error(err)
			}

			t.Logf("EN story: %+v", storyEN)
		})
	}
}

func TestContentEnv_RelatedStories(t *testing.T) {
	env := NewContentEnv(DB)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve related stories",
			args:    args{id: "001076306"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RelatedStories(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelatedStories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Related stories: %+v", got)
		})
	}
}
