package util

import (
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	type args struct {
		m map[string]bool
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single key",
			args: args{
				m: map[string]bool{
					"key1": true,
				},
			},
			want: []string{"key1"},
		},
		{
			name: "multiple keys",
			args: args{
				m: map[string]bool{
					"key1": true,
					"key2": true,
					"key3": true,
				},
			},
			want: []string{"key1", "key2", "key3"},
		},
		{
			name: "no keys",
			args: args{
				m: map[string]bool{},
			},
			want: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Keys(test.args.m)
			if len(got) != len(test.want) {
				t.Errorf("Keys() got length %v want %v", len(got), len(test.want))
			}
			sort.Strings(got) //order doesn't matter, sort to confirm values are correct
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("Keys()[%v] got %v want %v", i, got[i], test.want[i])
				}
			}
		})

	}

}
