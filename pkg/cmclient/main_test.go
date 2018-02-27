package cmclient

import "testing"

func TestUploadToServer(t *testing.T) {
	type args struct {
		filePaths      []string
		serverEndpoint string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "false test",
			args: args{
				filePaths:      []string{},
				serverEndpoint: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadToServer(tt.args.filePaths, tt.args.serverEndpoint); (err != nil) != tt.wantErr {
				t.Errorf("UploadToServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
