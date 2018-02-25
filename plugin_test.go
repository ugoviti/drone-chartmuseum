package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_getUnique(t *testing.T) {
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
			if got := getUnique(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDir(t *testing.T) {
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
			if got := isDir(tt.args.filePath); got != tt.want {
				t.Errorf("isDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUniqueParentFolders(t *testing.T) {
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
			if got := getUniqueParentFolders(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUniqueParentFolders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterExtFiles(t *testing.T) {
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
			trans := cmp.Transformer("Sort", sortStringSlice)
			if got := filterExtFiles(tt.args.files); cmp.Equal(got, tt.want, trans) {
				t.Errorf("filterExtFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDiffFiles(t *testing.T) {
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
			if got := getDiffFiles(tt.args.repoPath, tt.args.previousCommitID, tt.args.commitID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDiffFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveChartToPackage(t *testing.T) {
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
			got, err := saveChartToPackage(tt.args.chartPath, tt.args.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("saveChartToPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("saveChartToPackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uploadToServer(t *testing.T) {
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
			uploadToServer(tt.args.filePaths, tt.args.serverEndpoint)
		})
	}
}
