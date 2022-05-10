package utils

import (
	"reflect"
	"testing"

	"birthday-greeting/types"
)

var (
	GreetingPictureURL = "somewhere from S3"
)

func TestCraftBirthdayGreetingForUser(t *testing.T) {
	type args struct {
		user               types.User
		greetingPictureURL string
	}
	tests := []struct {
		name    string
		args    args
		want    types.BirthdayGreeting
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				user:               types.User{FirstName: "John", Gender: "M"},
				greetingPictureURL: GreetingPictureURL,
			},
			want: types.BirthdayGreeting{
				Title:   GreetingTitle,
				Content: "Happy birthday, dear John!",
				Picture: GreetingPictureURL,
			},
			wantErr: false,
		},
		{
			name: "Invalid FirstName",
			args: args{
				user:               types.User{FirstName: "", Gender: "F"},
				greetingPictureURL: "somewhere from S3",
			},
			want:    types.BirthdayGreeting{},
			wantErr: true,
		},
		{
			name: "Invalid URL",
			args: args{
				user:               types.User{FirstName: "Ellen", Gender: "T"},
				greetingPictureURL: "",
			},
			want:    types.BirthdayGreeting{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := CraftBirthdayGreetingForUser(tt.args.user, tt.args.greetingPictureURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CraftBirthdayGreetingForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
