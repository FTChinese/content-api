package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/content-api/pkg/db"
	"go.uber.org/zap/zaptest"
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

	DB = db.MockDB()
	if err != nil {
		panic(err)
	}
}

func TestStoryEnv_retrieveRawStory(t *testing.T) {

	env := StoryEnv{db: DB}

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
			got, err := env.retrieveRawStory(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveRawStory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Raw story: %+v\n", got)

			storyEN, err := got.BuildEN()
			if err != nil {
				t.Error(err)
			}

			t.Logf("EN story: %+v", storyEN)
		})
	}
}

func TestStoryEnv_retrieveRelatedStories(t *testing.T) {
	env := NewStoryEnv(DB, zaptest.NewLogger(t))

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

			got, err := env.retrieveRelatedStories(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveRelatedStories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Related stories: %+v", got)
		})
	}
}
