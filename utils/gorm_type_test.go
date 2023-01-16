package utils

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
)

func TestUUIDBinary_Value(t *testing.T) {

	tests := []struct {
		name    string
		id      UUIDBinary
		wantV   driver.Value
		wantErr bool
	}{
		{
			name:    "ok",
			id:      "cbec9655-6742-4398-b9e5-bc16e9a19843",
			wantV:   []uint8{203, 236, 150, 85, 103, 66, 67, 152, 185, 229, 188, 22, 233, 161, 152, 67},
			wantErr: false,
		},
		{
			name:    "err string nil",
			id:      "",
			wantV:   nil,
			wantErr: false,
		},
		{
			name:    "uuid.Parse err",
			id:      "123",
			wantErr: true,
		},
	}
	for k, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("case:", k)
			gotV, err := tt.id.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUIDBinary.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotV, tt.wantV) {
				fmt.Println(reflect.TypeOf(gotV))
				t.Errorf("UUIDBinary.Value() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

var TestUuid UUIDBinary = "cbec9655-6742-4398-b9e5-bc16e9a19843"

func TestUUIDBinary_String(t *testing.T) {
	tests := []struct {
		name string
		id   UUIDBinary
		want string
	}{
		{
			name: "ok",
			id:   "cbec9655-6742-4398-b9e5-bc16e9a19843",
			want: "cbec9655-6742-4398-b9e5-bc16e9a19843",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("UUIDBinary.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDBinary_Binary(t *testing.T) {
	tests := []struct {
		name    string
		id      UUIDBinary
		want    []byte
		wantErr bool
	}{
		{
			name:    "ok",
			id:      "cbec9655-6742-4398-b9e5-bc16e9a19843",
			want:    []byte{203, 236, 150, 85, 103, 66, 67, 152, 185, 229, 188, 22, 233, 161, 152, 67},
			wantErr: false,
		},
		{
			name:    "error",
			id:      "123",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.Binary()
			if (err != nil) != tt.wantErr {
				t.Errorf("UUIDBinary.Binary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUIDBinary.Binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDBinary_MustBinary(t *testing.T) {
	tests := []struct {
		name string
		id   UUIDBinary
		want []byte
	}{
		{
			name: "ok",
			id:   "cbec9655-6742-4398-b9e5-bc16e9a19843",
			want: []byte{203, 236, 150, 85, 103, 66, 67, 152, 185, 229, 188, 22, 233, 161, 152, 67},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.MustBinary(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UUIDBinary.MustBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUIDBinary_Scan(t *testing.T) {
	type args struct {
		value interface{}
	}

	var a1 args
	a1.value = ""
	var u1 UUIDBinary = "cbec9655-6742-4398-b9e5-bc16e9a19843"

	tests := []struct {
		name    string
		id      *UUIDBinary
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			id:      &u1,
			args:    a1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("UUIDBinary.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringArray_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	var array1 StringArray

	tests := []struct {
		name    string
		a       *StringArray
		args    args
		wantErr bool
	}{

		{
			name: "error",
			a:    &array1,
			args: args{
				src: 1,
			},
			wantErr: true,
		},
	}
	for K, tt := range tests {
		fmt.Println("test case :", K)
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("StringArray.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringArray_Value(t *testing.T) {
	var array1 StringArray
	array1 = append(array1, "1", "2", "3")
	var array2 StringArray

	tests := []struct {
		name    string
		a       StringArray
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "ok",
			a:       array1,
			want:    `["1","2","3"]`,
			wantErr: false,
		},
		{
			name:    "StringArray is nil",
			a:       array2,
			want:    "[]",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("StringArray.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringArray.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
