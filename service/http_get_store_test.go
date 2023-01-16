package service

import (
	"reflect"
	"testing"
)

func TestGetStore(t *testing.T) {
	type args struct {
		StoreId string
	}
	tests := []struct {
		name    string
		args    args
		want    *Store
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				StoreId: "123",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStore(tt.args.StoreId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStoreFromTotoro(t *testing.T) {
	type args struct {
		StoreId string
	}
	tests := []struct {
		name    string
		args    args
		want    *StoreInfo
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				StoreId: "123",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStoreFromTotoro(tt.args.StoreId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStoreFromTotoro() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStoreFromTotoro() = %v, want %v", got, tt.want)
			}
		})
	}
}
