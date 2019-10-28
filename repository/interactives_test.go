package repository

import (
	"github.com/FTChinese/go-rest"
	"testing"
)

func TestInteractiveEnv_RetrieveChannelTeasers(t *testing.T) {
	env := NewInteractiveEnv(DB)

	type args struct {
		keyWords string
		p        gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "麦可林学英语",
			args: args{
				keyWords: "麦可林学英语",
				p:        gorest.NewPagination(1, 12),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveChannelTeasers(tt.args.keyWords, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveChannelTeasers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestInteractiveEnv_RetrieveAudioArticle(t *testing.T) {
	env := NewInteractiveEnv(DB)
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve article for audio",
			args: args{
				id: 14168,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveAudioArticle(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveAudioArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestInteractiveEnv_RetrieveSpeedReadingTeasers(t *testing.T) {
	env := NewInteractiveEnv(DB)

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve speed reading teasers",
			args: args{
				p: gorest.NewPagination(1, 12),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveSpeedReadingTeasers(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpeedReadingTeasers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Speed reading teasers: %+v", got)
		})
	}
}

func TestInteractiveEnv_RetrieveSpeedReading(t *testing.T) {
	env := NewInteractiveEnv(DB)

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Article for speed ing reading",
			args:    args{id: 14156},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.RetrieveSpeedReading(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveSpeedReading() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("An article for speed reading: %+v", got)
		})
	}
}
