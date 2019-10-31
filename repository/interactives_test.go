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

			got, err := env.RetrieveTeasers(tt.args.keyWords, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveTeasers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}

func TestInteractiveEnv_RetrieveRawContent(t *testing.T) {
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
			name:    "FTA Report",
			args:    args{id: 14065},
			wantErr: false,
		},
		{
			name:    "Speed Reading",
			args:    args{id: 14058},
			wantErr: false,
		},
		{
			name:    "FT Radio",
			args:    args{id: 14172},
			wantErr: false,
		},
		{
			name:    "MLE",
			args:    args{id: 14168},
			wantErr: false,
		},
		{
			name:    "Daily Word",
			args:    args{id: 14162},
			wantErr: false,
		},
		{
			name:    "Book",
			args:    args{id: 14148},
			wantErr: false,
		},
		{
			name:    "Boomear",
			args:    args{id: 14116},
			wantErr: false,
		},
		{
			name:    "Life of a song",
			args:    args{id: 7825},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := env.RetrieveRawContent(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveRawContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
