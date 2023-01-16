package models

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
)

func TestGetAppTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock.New() err: %s", err)
	}
	defer db.Close()
	Gdb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("gorm.Open() err: %s", err)
	}
	Gdb.LogMode(true)

	var tag1 = AppTag{
		Id:            uint64(1),
		ApplicationId: uint64(1),
		TagId:         uint64(1),
	}

	// mock.ExpectQuery("^SELECT * FROM `oauth_applications`  WHERE id= 1 LIMIT 1").
	mock.ExpectQuery("^SELECT *").WithArgs(tag1.ApplicationId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "tag_id"}).
			AddRow(tag1.Id, tag1.ApplicationId, tag1.TagId))

	type args struct {
		db *gorm.DB
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []AppTag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db: Gdb,
				id: tag1.ApplicationId,
			},
			want: []AppTag{
				tag1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppTags(tt.args.db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock.New() err: %s", err)
	}
	defer db.Close()
	Gdb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("gorm.Open() err: %s", err)
	}
	Gdb.LogMode(true)

	var tag1 = Tag{
		Id:       uint64(1),
		NameZh:   "namezh1",
		TagLevel: "taglevel1",
		ParentId: uint64(1),
		NameEn:   "nameen1",
	}

	// mock.ExpectQuery("^SELECT * FROM `oauth_applications`  WHERE id= 1 LIMIT 1").
	mock.ExpectQuery("^SELECT *").WithArgs(tag1.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name_zh", "tag_level", "parent_id", "name_en"}).
			AddRow(tag1.Id, tag1.NameZh, tag1.TagLevel, tag1.ParentId, tag1.NameEn))

	type args struct {
		db *gorm.DB
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Tag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db: Gdb,
				id: tag1.Id,
			},
			want:    &tag1,
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				db: Gdb,
				id: tag1.Id,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTag(tt.args.db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteAppTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("sqlmock.New() err: %s", err)
	}
	defer db.Close()
	Gdb, err := gorm.Open("mysql", db)
	if err != nil {
		t.Errorf("gorm.Open() err: %s", err)
	}
	Gdb.LogMode(true)
	apptag := new(AppTag)
	apptag.Id = uint64(1)
	mock.ExpectExec(regexp.QuoteMeta("DELETE")).WithArgs(apptag.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	type args struct {
		db *gorm.DB
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db: Gdb,
				id: apptag.Id,
			},
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				db: Gdb,
				id: uint64(2),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteAppTags(tt.args.db, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAppTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertTagsToTagViews(t *testing.T) {
	type args struct {
		tags []Tag
	}

	a1 := args{
		tags: []Tag{
			{Id: 1, NameZh: "1", NameEn: "1", ParentId: 0, TagLevel: "1"},
			{Id: 2, NameZh: "2", NameEn: "2", ParentId: 0, TagLevel: "1"},
			{Id: 3, NameZh: "3", NameEn: "3", ParentId: 0, TagLevel: "1"},
			{Id: 4, NameZh: "4", NameEn: "4", ParentId: 0, TagLevel: "1"},
			{Id: 5, NameZh: "5", NameEn: "5", ParentId: 1, TagLevel: "2"},
			{Id: 6, NameZh: "6", NameEn: "6", ParentId: 2, TagLevel: "2"},
			{Id: 7, NameZh: "7", NameEn: "7", ParentId: 3, TagLevel: "2"},
			{Id: 8, NameZh: "8", NameEn: "8", ParentId: 4, TagLevel: "2"},
			{Id: 9, NameZh: "9", NameEn: "9", ParentId: 1, TagLevel: "2"},
			{Id: 10, NameZh: "10", NameEn: "10", ParentId: 2, TagLevel: "2"},
			{Id: 11, NameZh: "11", NameEn: "11", ParentId: 3, TagLevel: "2"},
			{Id: 12, NameZh: "12", NameEn: "12", ParentId: 4, TagLevel: "2"},
			{Id: 13, NameZh: "13", NameEn: "13", ParentId: 12, TagLevel: "3"},
			{Id: 14, NameZh: "14", NameEn: "14", ParentId: 5, TagLevel: "3"},
			{Id: 15, NameZh: "15", NameEn: "15", ParentId: 9, TagLevel: "3"},
		},
	}
	tests := []struct {
		name    string
		args    args
		want    []TagView
		wantErr bool
	}{
		{
			name: "成功",
			args: a1,
			want: []TagView{
				{1, "1", "1", []TagView{
					{5, "5", "5", []TagView{
						{14, "14", "14", nil},
					}},
					{9, "9", "9", []TagView{
						{15, "15", "15", nil},
					}},
				}},
				{2, "2", "2", []TagView{
					{6, "6", "6", nil},
					{10, "10", "10", nil},
				}},
				{3, "3", "3", []TagView{
					{7, "7", "7", nil},
					{11, "11", "11", nil},
				}},
				{4, "4", "4", []TagView{
					{8, "8", "8", nil},
					{12, "12", "12", []TagView{
						{13, "13", "13", nil},
					}},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTagsToTagViews(tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTagsToTagViews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println("got:", reflect.TypeOf(got), got)
				fmt.Println("want:", reflect.TypeOf(tt.want), tt.want)

				t.Errorf("ConvertTagsToTagViews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateAppTag(t *testing.T) {
	type args struct {
		db     *gorm.DB
		appTag AppTag
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:     config.DB,
				appTag: AppTag{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateAppTag(tt.args.db, tt.args.appTag); (err != nil) != tt.wantErr {
				t.Errorf("CreateAppTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []Tag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				config.DB,
			},
			want:    make([]Tag, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTags(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTagsByParentId(t *testing.T) {
	type args struct {
		db       *gorm.DB
		parentId []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []Tag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:       config.DB,
				parentId: []uint64{123},
			},
			want:    make([]Tag, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTagsByParentId(tt.args.db, tt.args.parentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsByParentId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsByParentId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTagsByLevel(t *testing.T) {
	type args struct {
		db    *gorm.DB
		level string
	}
	tests := []struct {
		name    string
		args    args
		want    []Tag
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				config.DB,
				"1",
			},
			want:    make([]Tag, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTagsByLevel(tt.args.db, tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsByLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsByLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppCountGroubyTag(t *testing.T) {
	type args struct {
		db     *gorm.DB
		appIds []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []TagInfo
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:     config.DB,
				appIds: []uint64{1},
			},
			want:    make([]TagInfo, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppCountGroubyTag(tt.args.db, tt.args.appIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppCountGroubyTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppCountGroubyTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublishedAppByTagId(t *testing.T) {
	type args struct {
		db     *gorm.DB
		tagIds []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "ok1",
			args: args{
				db:     config.DB,
				tagIds: []uint64{1},
			},
			want:    make([]uint64, 0),
			wantErr: false,
		},
		{
			name: "ok2",
			args: args{
				db:     config.DB,
				tagIds: []uint64{},
			},
			want:    make([]uint64, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublishedAppByTagId(tt.args.db, tt.args.tagIds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublishedAppByTagId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.DeepEqual(got, tt.want) {
				fmt.Println(reflect.TypeOf(got), got)
				fmt.Println(reflect.TypeOf(tt.want), tt.want)

			}
		})
	}
}

func TestGerAppIdsOrderbyInstalledCount(t *testing.T) {
	type args struct {
		db     *gorm.DB
		appIds []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:     config.DB,
				appIds: []uint64{1},
			},
			want:    []uint64{1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GerAppIdsOrderbyInstalledCount(tt.args.db, tt.args.appIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GerAppIdsOrderbyInstalledCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GerAppIdsOrderbyInstalledCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppTagsByIds(t *testing.T) {
	type args struct {
		db  *gorm.DB
		ids []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []TagView
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:  config.DB,
				ids: []uint64{1},
			},
			want:    make([]TagView, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppTagsByIds(tt.args.db, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppTagsByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppTagsByIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppsTagsByIds(t *testing.T) {
	type args struct {
		db  *gorm.DB
		ids []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []TagAppView
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				db:  config.DB,
				ids: []uint64{1},
			},
			want:    make([]TagAppView, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppsTagsByIds(tt.args.db, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppsTagsByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAppsTagsByIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagAppViewSliceGroupByAppid(t *testing.T) {
	type args struct {
		tagAppViews []TagAppView
	}
	tests := []struct {
		name string
		args args
		want map[uint64][]TagAppView
	}{
		{
			name: "ok",
			args: args{
				tagAppViews: []TagAppView{
					{
						AppId:  1,
						TagId:  1,
						NameZh: "z",
						NameEn: "e",
					},
				},
			},
			want: map[uint64][]TagAppView{
				1: {
					{
						AppId:  1,
						TagId:  1,
						NameZh: "z",
						NameEn: "e",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TagAppViewSliceGroupByAppid(tt.args.tagAppViews); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TagAppViewSliceGroupByAppid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAppTagViewsParent(t *testing.T) {
	type args struct {
		tavs     []TagAppView
		parentId uint64
	}

	tests := []struct {
		name string
		args args
		want *TagAppView
	}{
		{
			name: "ok",
			args: args{
				tavs: []TagAppView{
					{TagId: 1},
					{TagId: 2},
					{TagId: 3},
					{TagId: 4},
				},
				parentId: 4,
			},
			want: &TagAppView{TagId: 4},
		},
		{
			name: "ok1",
			args: args{
				tavs: []TagAppView{
					{TagId: 1},
					{TagId: 2},
					{TagId: 3},
					{TagId: 4},
				},
				parentId: 1,
			},
			want: &TagAppView{TagId: 1},
		},
		{
			name: "ok1",
			args: args{
				tavs: []TagAppView{
					{TagId: 1},
					{TagId: 2},
					{TagId: 3},
					{TagId: 4},
				},
				parentId: 5,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAppTagViewsParent(tt.args.tavs, tt.args.parentId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAppTagViewsParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimAppTagViews(t *testing.T) {
	type args struct {
		tav []TagAppView
	}
	tests := []struct {
		name string
		args args
		want []TagAppView
	}{
		{
			name: "ok",
			args: args{
				tav: []TagAppView{
					{NameZh: "其他", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他1", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
					{NameZh: "一级", ParentId: 1, TagId: 1, TagLevel: "1"},
				},
			},
			want: []TagAppView{
				{NameZh: "一级", ParentId: 1, TagId: 1, TagLevel: "1"},
				{NameZh: "其他1", ParentId: 1, TagLevel: "2"},
				{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
				{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
			},
		},
		{
			name: "ok1",
			args: args{
				tav: []TagAppView{
					{NameZh: "其他1", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
					{NameZh: "一级", ParentId: 1, TagId: 1, TagLevel: "1"},
				},
			},
			want: []TagAppView{
				{NameZh: "其他1", ParentId: 1, TagLevel: "2"},
				{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
				{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
			},
		},
		{
			name: "ok2",
			args: args{
				tav: []TagAppView{
					{NameZh: "其他1", ParentId: 5, TagLevel: "2"},
					{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
					{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
					{NameZh: "一级", ParentId: 1, TagId: 1, TagLevel: "1"},
				},
			},
			want: []TagAppView{
				{NameZh: "其他1", ParentId: 5, TagLevel: "2"},
				{NameZh: "其他2", ParentId: 1, TagLevel: "2"},
				{NameZh: "其他3", ParentId: 1, TagLevel: "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimAppTagViews(tt.args.tav); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimAppTagViews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublishedAppByTagIdAndAppRegion(t *testing.T) {
	type args struct {
		db     *gorm.DB
		region string
		tagIds []uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "ok1",
			args: args{
				db:     config.DB,
				region: "CN",
				tagIds: []uint64{1},
			},
			want:    []uint64{},
			wantErr: false,
		},
		{
			name: "ok2",
			args: args{
				db:     config.DB,
				region: "Global",
				tagIds: []uint64{},
			},
			want:    []uint64{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPublishedAppByTagIdAndAppRegion(tt.args.db, tt.args.region, tt.args.tagIds...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublishedAppByTagIdAndAppRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.want) {
				fmt.Println(reflect.TypeOf(got), got)
				fmt.Println(reflect.TypeOf(tt.want), tt.want)

			}
		})
	}
}
