package main

import "testing"

func TestPlugin_defaultExec(t *testing.T) {
	type fields struct {
		Config Config
	}
	type args struct {
		files []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Config: tt.fields.Config,
			}
			p.defaultExec(tt.args.files)
		})
	}
}

func TestPlugin_exec(t *testing.T) {
	type fields struct {
		Config Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{
				Config: tt.fields.Config,
			}
			if err := p.exec(); (err != nil) != tt.wantErr {
				t.Errorf("Plugin.exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
