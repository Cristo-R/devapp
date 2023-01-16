package utils

import "testing"

func TestSlicePage(t *testing.T) {
	type args struct {
		page     int
		pageSize int
		nums     int
	}
	tests := []struct {
		name           string
		args           args
		wantSliceStart int
		wantSliceEnd   int
	}{
		{
			name: "ok1",
			args: args{
				page:     -1,
				pageSize: -1,
				nums:     1,
			},
			wantSliceStart: 0,
			wantSliceEnd:   1,
		},
		{
			name: "ok2",
			args: args{
				page:     2,
				pageSize: 1,
				nums:     1,
			},
			wantSliceStart: 0,
			wantSliceEnd:   0,
		},
		{
			name: "ok3",
			args: args{
				page:     1,
				pageSize: 1,
				nums:     1,
			},
			wantSliceStart: 0,
			wantSliceEnd:   1,
		},
		{
			name: "ok4",
			args: args{
				page:     2,
				pageSize: 2,
				nums:     4,
			},
			wantSliceStart: 2,
			wantSliceEnd:   4,
		},
		{
			name: "error",
			args: args{
				page:     5,
				pageSize: 5,
				nums:     -23,
			},
			wantSliceStart: 0,
			wantSliceEnd:   -23,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSliceStart, gotSliceEnd := SlicePage(tt.args.page, tt.args.pageSize, tt.args.nums)
			if gotSliceStart != tt.wantSliceStart {
				t.Errorf("SlicePage() gotSliceStart = %v, want %v", gotSliceStart, tt.wantSliceStart)
			}
			if gotSliceEnd != tt.wantSliceEnd {
				t.Errorf("SlicePage() gotSliceEnd = %v, want %v", gotSliceEnd, tt.wantSliceEnd)
			}
		})
	}
}
