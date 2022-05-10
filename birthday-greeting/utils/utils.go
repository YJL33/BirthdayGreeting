package utils

import (
	"encoding/json"
	"errors"

	"birthday-greeting/types"
)

var (
	GreetingTitle = "Subject: Happy birthday!"
)

func CraftBirthdayGreetingForUser(user types.User) (types.BirthdayGreeting, error) {
	if user.FirstName == "" || user.LastName == "" {
		userInJson, _ := json.Marshal(user)
		return types.BirthdayGreeting{}, errors.New("FirstName or LastName is missing:" + string(userInJson))
	}

	greetingContent := "Happy birthday, dear " + user.LastName + ", " + user.FirstName + "!"
	return types.BirthdayGreeting{
		Title:   GreetingTitle,
		Content: greetingContent,
	}, nil
}
