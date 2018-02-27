package main

import (
	"reflect"
	"testing"
)

func TestPlugin_GetDiffFiles(t *testing.T) {
	pTest := &Plugin{
		Config: Config{
			PreviousCommitID: "12345",
			CurrentCommitID:  "67890",
		},
	}
	tests := []struct {
		name    string
		p       *Plugin
		want    []string
		wantErr bool
	}{
		{
			name: "true test",
			p:    pTest,
			// want:    []string{"123"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetDiffFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.GetDiffFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.GetDiffFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlugin_SaveChartToPackage(t *testing.T) {
	pTest := &Plugin{
		Config: Config{
			SaveDir: "test/uploads/",
		},
	}
	type args struct {
		chartPath string
	}
	tests := []struct {
		name        string
		p           *Plugin
		args        args
		wantMessage string
		wantErr     bool
	}{
		{
			name: "true test",
			p:    pTest,
			args: args{
				chartPath: "test/chart/",
			},
			wantMessage: "test/uploads/chart-0.1.0.tgz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := tt.p.SaveChartToPackage(tt.args.chartPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Plugin.SaveChartToPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("Plugin.SaveChartToPackage() = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}

func TestPlugin_exec(t *testing.T) {
	tests := []struct {
		name    string
		p       *Plugin
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.exec(); (err != nil) != tt.wantErr {
				t.Errorf("Plugin.exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlugin_PackageAndUpload(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		p       *Plugin
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.PackageAndUpload(tt.args.files); (err != nil) != tt.wantErr {
				t.Errorf("Plugin.PackageAndUpload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
