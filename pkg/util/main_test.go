package util

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_GetUnique(t *testing.T) {
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

func Test_IsDir(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
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
			if got := IsDir(tt.args.filePath); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetParentFolders(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "expect to be false",
			args: args{
				files: []string{"test1/file.txt", "test2/file.txt"},
			},
			want: []string{},
		},
		{
			name: "expect to be true",
			args: args{
				files: []string{"test/file.txt", "test/folder/"},
			},
			want: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetParentFolders(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParentFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FilterExtFiles(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "true",
			args: args{
				files: []string{"test/file.txt", "test/file.yml", "test/file.yaml"},
			},
			want: []string{"test/file.yaml", "test/file.yml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trans := cmp.Transformer("Sort", SortStringSlice)
			if got := FilterExtFiles(tt.args.files); cmp.Equal(got, tt.want, trans) {
				t.Errorf("FilterExtFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetDiffFiles(t *testing.T) {
	type args struct {
		repoPath         string
		previousCommitID string
		commitID         string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDiffFiles(tt.args.repoPath, tt.args.previousCommitID, tt.args.commitID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDiffFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SaveChartToPackage(t *testing.T) {
	type args struct {
		chartPath string
		dstPath   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SaveChartToPackage(tt.args.chartPath, tt.args.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveChartToPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SaveChartToPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UploadToServer(t *testing.T) {
	type args struct {
		filePaths      []string
		serverEndpoint string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UploadToServer(tt.args.filePaths, tt.args.serverEndpoint)
		})
	}
}
