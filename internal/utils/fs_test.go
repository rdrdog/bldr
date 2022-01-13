package utils

import (
	"testing"

	"github.com/spf13/afero"
)

func TestFileExists(t *testing.T) {

	var AppFs = afero.NewOsFs()

	AppFs.Create("test.txt")

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "If file exists",
			want:    true,
			wantErr: false,
			args:    args{name: "test.txt"},
		},
		{
			name:    "If file doesn't exist",
			want:    false,
			wantErr: false,
			args:    args{name: "../test_false.txt"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FileExists(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
	AppFs.Remove("test.txt")
}
