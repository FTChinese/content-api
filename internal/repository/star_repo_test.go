package repository

import (
	"reflect"
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/content-api/internal/pkg"
	"gitlab.com/ftchinese/content-api/pkg/db"
	"gitlab.com/ftchinese/content-api/pkg/ids"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestStarRepo_StarArticle(t *testing.T) {

	env := NewStarRepo(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		params pkg.StarParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "start article",
			args: args{
				params: pkg.StarParams{
					UserID: "1234567890",
					ID:     "001076303",
					Kind:   pkg.ContentKindStory,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.StarArticle(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("StarRepo.StarArticle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStarRepo_listStarred(t *testing.T) {
	type fields struct {
		dbs    db.ReadWriteMyDBs
		logger *zap.Logger
	}
	type args struct {
		id ids.UserIDs
		p  gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []pkg.RawStarred
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StarRepo{
				dbs:    tt.fields.dbs,
				logger: tt.fields.logger,
			}
			got, err := env.listStarred(tt.args.id, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("StarRepo.listStarred() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StarRepo.listStarred() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStarRepo_ListStarred(t *testing.T) {
	env := NewStarRepo(db.MockMySQL(), zaptest.NewLogger(t))

	type args struct {
		id ids.UserIDs
		p  gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    pkg.PagedList[pkg.Starred]
		wantErr bool
	}{
		{
			name: "list articles",
			args: args{
				id: ids.UserIDs{
					FtcID: null.StringFrom("1234567890"),
				},
				p: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListStarred(tt.args.id, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("StarRepo.ListStarred() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("StarRepo.ListStarred() = %v, want %v", got, tt.want)
			// }

			t.Logf("%v", got)
		})
	}
}
