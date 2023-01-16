package utils

import (
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func Test_getRequestParamStr(t *testing.T) {
	type args struct {
		data map[string]interface{}
	}

	var array_test []interface{}
	array_test = append(array_test, 1, 2, 3)
	var data1 = map[string]interface{}{
		"1": 1,
		"2": 2,
	}
	var data2 = map[string]interface{}{
		"1": 1,
		"2": array_test,
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok:string, float32, float64, int, int8, int16, int32, int64 ",
			args: args{data: data1},
			want: "1=1&2=2",
		},
		{
			name: "ok:[]interface{}",
			args: args{data: data2},
			want: "1=1&2[]=1&2[]=2&2[]=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRequestParamStr(tt.args.data); got != tt.want {
				t.Errorf("getRequestParamStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		data   interface{}
	}

	api := "http://1sdnfjkl.com"
	mockResponse := "mock response body"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", api, httpmock.NewStringResponder(200, string(mockResponse)))

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 int
	}{
		{
			name: "get data nil ",
			args: args{
				method: "GET",
				url:    "",
				data:   nil,
			},
			want:  nil,
			want1: 500,
		},
		{
			name: "post datat nil",
			args: args{
				method: "POST",
				url:    "",
				data:   nil,
			},
			want:  nil,
			want1: 500,
		},
		{
			name: "post datat nil",
			args: args{
				method: "POST",
				url:    api,
				data:   nil,
			},
			want:  nil,
			want1: 500,
		},
		// {
		// 	name: "post datat nil",
		// 	args: args{
		// 		method: "GET",
		// 		url:    api1,
		// 		data:   nil,
		// 	},
		// 	want:  nil,
		// 	want1: 500,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := HttpRequest(tt.args.method, tt.args.url, tt.args.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpRequest() got = %v, want %v", string(got), string(tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("HttpRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
