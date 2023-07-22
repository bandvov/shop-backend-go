package main

import (
	"testing"
)

func Test_getEnvVariable(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test get env variables",
			args: args{
				name: "test",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEnvVariable(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEnvVariable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEnvVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}
