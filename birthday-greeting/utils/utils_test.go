package utils

import (
	"reflect"
	"testing"

	"birthday-greeting/types"
)

func TestCraftBirthdayGreetingForUser(t *testing.T) {
	type args struct {
		user types.User
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
				user: types.User{
					FirstName: "John",
					Gender:    "M",
				},
			},
			want: types.BirthdayGreeting{
				Title:   GreetingTitle,
				Content: "Happy birthday, dear John!" + MessageForSpecialDiscount + DiscountItemForMale,
			},
			wantErr: false,
		},
		{
			name: "Invalid FirstName",
			args: args{
				user: types.User{
					FirstName: "",
					Gender:    "F",
				},
			},
			want:    types.BirthdayGreeting{},
			wantErr: true,
		},
		{
			name: "Non-binary Gender",
			args: args{
				user: types.User{
					FirstName: "Ellen",
					Gender:    "T",
				},
			},
			want:    types.BirthdayGreeting{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := CraftBirthdayGreetingForUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CraftBirthdayGreetingForUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
