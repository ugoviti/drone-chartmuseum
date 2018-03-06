package util

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false test",
			args: args{
				s: []string{"test1", "test2"},
				e: "test3",
			},
			want: false,
		},
		{
			name: "true test",
			args: args{
				s: []string{"test1", "test2"},
				e: "test1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteEmpty(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name  string
		args  args
		wantR []string
	}{
		{
			name: "true test",
			args: args{
				s: []string{"", "test1"},
			},
			wantR: []string{"test1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := DeleteEmpty(tt.args.s); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("DeleteEmpty() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestExtractDirs(t *testing.T) {
	type args struct {
		fileInfos []os.FileInfo
	}
	fileInfos, _ := ioutil.ReadDir("test/")
	tests := []struct {
		name           string
		args           args
		wantResultList []string
	}{
		{
			name: "true test",
			args: args{
				fileInfos: fileInfos,
			},
			wantResultList: []string{"file.txt", "file.yaml", "file.yml", "folder"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResultList := ExtractDirs(tt.args.fileInfos); !reflect.DeepEqual(gotResultList, tt.wantResultList) {
				t.Errorf("ExtractDirs() = %v, want %v", gotResultList, tt.wantResultList)
			}
		})
	}
}

func TestSortStringSlice(t *testing.T) {
	type args struct {
		in []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "true test",
			args: args{
				in: []string{"a", "c", "b"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortStringSlice(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUnique(t *testing.T) {
	type args struct {
		input []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "expect to be true",
			args: args{
				input: []string{"component1", "component1", "component2", "component3"},
			},
			want: []string{"component1", "component2", "component3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUnique(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "folder",
			args: args{
				filePath: "test/folder",
			},
			want: true,
		},
		{
			name: "file",
			args: args{
				filePath: "test/file.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDir(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetParentFolders(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name            string
		args            args
		wantResultSlice []string
	}{
		{
			name: "expect to be false",
			args: args{
				files: []string{"test1/file.txt", "test2/file.txt"},
			},
			wantResultSlice: []string{},
		},
		{
			name: "expect to be true",
			args: args{
				files: []string{"test/file.txt", "test/folder/"},
			},
			wantResultSlice: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResultSlice := GetParentFolders(tt.args.files); !reflect.DeepEqual(gotResultSlice, tt.wantResultSlice) {
				t.Errorf("GetParentFolders() = %v, want %v", gotResultSlice, tt.wantResultSlice)
			}
		})
	}
}

func TestFilterExtFiles(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name            string
		args            args
		wantResultSlice []string
	}{
		{
			name: "true",
			args: args{
				files: []string{"test/file.txt", "test/file.yml", "test/file.yaml"},
			},
			wantResultSlice: []string{"test/file.yaml", "test/file.yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResultSlice := FilterExtFiles(tt.args.files); !reflect.DeepEqual(gotResultSlice, tt.wantResultSlice) {
				t.Errorf("FilterExtFiles() = %v, want %v", gotResultSlice, tt.wantResultSlice)
			}
		})
	}
}
