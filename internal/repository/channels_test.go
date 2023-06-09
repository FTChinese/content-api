package repository

import (
	"reflect"
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTchinese/content-api/internal/pkg"
	"go.uber.org/zap/zaptest"
)

func TestChannelEnv_ListChannels(t *testing.T) {
	env := NewChannelEnv(DB, zaptest.NewLogger(t))

	got, err := env.ListChannels()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Channel list: %+v", got)

	t.Logf("Channel map: %+v", pkg.NewChannelMap(got))
}

func TestChannelEnv_retrieveChannel(t *testing.T) {
	env := NewChannelEnv(DB, zaptest.NewLogger(t))
	type args struct {
		pathName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "China channel",
			args: args{pathName: "china"},
		},
		{
			name: "China Economy Channel",
			args: args{pathName: "china_economy"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.retrieveChannel(tt.args.pathName)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("A channel: %+v", got)
		})
	}
}

func TestChannelEnv_RetrieveTeasers(t *testing.T) {

	env := NewChannelEnv(DB, zaptest.NewLogger(t))

	type args struct {
		channelID int64
		p         gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []pkg.Teaser
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveTeasers(tt.args.channelID, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveTeasers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveTeasers() got = %v, want %v", got, tt.want)
			}
		})
	}
}
