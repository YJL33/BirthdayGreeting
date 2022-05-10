package utils

import (
	"encoding/json"
	"errors"

	"birthday-greeting/types"
)

func CraftBirthdayGreetingForUser(user types.User) (types.BirthdayGreeting, error) {
	if user.FirstName == "" {
		userInJson, _ := json.Marshal(user)
		return types.BirthdayGreeting{}, errors.New("FirstName is missing:" + string(userInJson))
	}

	return types.BirthdayGreeting{
		Title:   "Subject: Happy birthday!",
		Content: "Happy birthday, dear " + user.FirstName + "!",
	}, nil
}
