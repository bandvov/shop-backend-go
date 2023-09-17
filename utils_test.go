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
	var errors = make(ValidationErrors)

	errors["email"] = []string{"email is required"}
	errors["password"] = []string{"password is required"}

	type args struct {
		body   interface{}
		fields []string
	}
	tests := []struct {
		name string
		args args
		want ValidationErrors
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
			want: make(ValidationErrors),
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
			want: errors,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validate(tt.args.body, tt.args.fields...)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() got1 = %v, want %v", got, tt.want)
			}
		})
	}
}
