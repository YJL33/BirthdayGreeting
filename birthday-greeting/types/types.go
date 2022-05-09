package types

// User is a struct representing user
type User struct {
	UserID      string `json:"ID"`
	DateOfBirth string `json:"DateOfBirth"`
	Email       string `json:"Email"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Gender      string `json:"Gender"`
}

type BirthdayGreeting struct {
	Title   string `json:"Title"`
	Content string `json:"Content"`
}
