package main

import (
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func Test_initApp(t *testing.T) {
	tests := []struct {
		name string
		want *cli.App
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initApp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := defaultAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("defaultAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
