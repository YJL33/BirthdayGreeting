package utils

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"birthday-greeting/types"
)

var (
	GreetingTitle = "Subject: Happy birthday!"
)

func CraftBirthdayGreetingForUser(user types.User, greetingPictureURL string) (types.BirthdayGreeting, error) {
	userInJson, _ := json.Marshal(user)
	if user.FirstName == "" {
		return types.BirthdayGreeting{}, errors.New("FirstName is missing:" + string(userInJson))
	}
	if greetingPictureURL == "" {
		return types.BirthdayGreeting{}, errors.New("PictureURL is missing:" + greetingPictureURL)
	}
	birthYear, err := strconv.Atoi(user.DateOfBirth[:4])
	if err != nil {
		return types.BirthdayGreeting{}, errors.New("Invalid birth year:" + string(userInJson))
	}
	currentYear, _ := strconv.Atoi(time.Now().Format("2006-01-02")[:4])
	if currentYear-birthYear <= 49 {
		return types.BirthdayGreeting{}, errors.New("User is too young:" + string(userInJson))
	}
	greetingContent := "Happy birthday, dear " + user.FirstName + "!"

	return types.BirthdayGreeting{
		Title:   GreetingTitle,
		Content: greetingContent,
		Picture: greetingPictureURL,
	}, nil
}
