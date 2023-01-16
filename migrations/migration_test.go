package migrations

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
)

func TestMigrate(t *testing.T) {
	if err := Migrate(config.DB); err != nil {
		panic(err)
	}

	Convey("check payment_extensions table", t, func() {
		So(config.DB.HasTable("payment_extensions"), ShouldBeTrue)
	})

	Convey("check account_alias table", t, func() {
		So(config.DB.HasTable("account_alias"), ShouldBeTrue)
	})
}

func TestIndexExist(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock.New() err: %s", err)
	}
	defer db.Close()
	Gdb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("gorm.Open() err: %s", err)
	}
	Gdb.LogMode(true)

	type args struct {
		db    *gorm.DB
		table string
		index string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:    Gdb,
				table: "1",
				index: "1",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, got := IndexExist(tt.args.db, tt.args.table, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
