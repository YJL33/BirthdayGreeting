package utils

import (
	"encoding/json"
	"errors"

	"birthday-greeting/types"
)

var (
	GreetingTitle = "Subject: Happy birthday!"
)

func CraftBirthdayGreetingForUser(user types.User, greetingPictureURL string) (types.BirthdayGreeting, error) {
	if user.FirstName == "" {
		userInJson, _ := json.Marshal(user)
		return types.BirthdayGreeting{}, errors.New("FirstName is missing:" + string(userInJson))
	}
	if greetingPictureURL == "" {
		return types.BirthdayGreeting{}, errors.New("PictureURL is missing:" + string(greetingPictureURL))
	}
	greetingContent := "Happy birthday, dear " + user.FirstName + "!"

	return types.BirthdayGreeting{
		Title:   GreetingTitle,
		Content: greetingContent,
		Picture: greetingPictureURL,
	}, nil
}
