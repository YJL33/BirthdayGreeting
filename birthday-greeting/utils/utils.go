package utils

import (
	"birthday-greeting/types"
)

func CraftBirthdayGreetingForUser(user types.User) types.BirthdayGreeting {
	return types.BirthdayGreeting{
		Title:   "Subject: Happy birthday!",
		Content: "Happy birthday, dear " + user.FirstName + "!",
	}
}
