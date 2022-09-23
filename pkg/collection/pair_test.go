package collection

import (
	"reflect"
	"testing"
)

func TestZipString(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []StringPair
	}{
		{
			name: "Equal length",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"1", "2", "3"},
			},
			want: []StringPair{
				{
					First:  "a",
					Second: "1",
				},
				{
					First:  "b",
					Second: "2",
				},
				{
					First:  "c",
					Second: "3",
				},
			},
		},
		{
			name: "Array a have more elements",
			args: args{
				a: []string{"a", "b", "c", "d"},
				b: []string{"1", "2", "3"},
			},
			want: []StringPair{
				{
					First:  "a",
					Second: "1",
				},
				{
					First:  "b",
					Second: "2",
				},
				{
					First:  "c",
					Second: "3",
				},
				{
					First:  "d",
					Second: "",
				},
			},
		},
		{
			name: "Array b have more elements",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"1", "2", "3", "4"},
			},
			want: []StringPair{
				{
					First:  "a",
					Second: "1",
				},
				{
					First:  "b",
					Second: "2",
				},
				{
					First:  "c",
					Second: "3",
				},
				{
					First:  "",
					Second: "4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZipString(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZipString() = %v, want %v", got, tt.want)
			}
		})
	}
}
