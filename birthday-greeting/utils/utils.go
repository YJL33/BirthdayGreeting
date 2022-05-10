package utils

import (
	"encoding/json"
	"errors"

	"birthday-greeting/types"
)

var (
	GreetingTitle             = "Subject: Happy birthday!"
	MessageForSpecialDiscount = "We offer special discount 20% off for the following items: \n"
	DiscountItemForMale       = "White Wine, iPhoneX"
	DiscountItemForFemale     = "Cosmetic, LV Handbags"
)

func CraftBirthdayGreetingForUser(user types.User) (types.BirthdayGreeting, error) {
	if user.FirstName == "" {
		userInJson, _ := json.Marshal(user)
		return types.BirthdayGreeting{}, errors.New("FirstName is missing:" + string(userInJson))
	}
	if user.Gender != "F" && user.Gender != "M" {
		userInJson, _ := json.Marshal(user)
		return types.BirthdayGreeting{}, errors.New("Gender is non-binary:" + string(userInJson))
	}
	var greetingContent string
	if user.Gender == "F" {
		greetingContent = "Happy birthday, dear " + user.FirstName + "!" + MessageForSpecialDiscount + DiscountItemForFemale
	} else {
		greetingContent = "Happy birthday, dear " + user.FirstName + "!" + MessageForSpecialDiscount + DiscountItemForMale
	}

	return types.BirthdayGreeting{
		Title:   GreetingTitle,
		Content: greetingContent,
	}, nil
}
