package main

import (
	"reflect"
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

func Test_validate(t *testing.T) {
	var errors = make(map[string][]string)

	errors["email"] = []string{"email is required"}
	errors["password"] = []string{"password is required"}

	type args struct {
		body   interface{}
		fields []string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 map[string][]string
	}{
		{
			name: "validate all required field valid",
			args: args{
				fields: []string{"email", "password"},
				body: User{
					Password: "23123",
					Email:    "test@test.aa",
				},
			},
			want:  false,
			want1: make(map[string][]string),
		},
		{
			name: "validate all required field not valid",
			args: args{
				fields: []string{"email", "password"},
				body: User{
					Password: "",
					Email:    "",
				},
			},
			want:  true,
			want1: errors,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := validate(tt.args.body, tt.args.fields...)
			if got != tt.want {
				t.Errorf("validate() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("validate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
