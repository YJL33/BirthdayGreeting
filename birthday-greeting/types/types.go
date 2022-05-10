package types

// User is a struct representing user
type User struct {
	UserID      string `json:"id"`
	DateOfBirth string `json:"dateOfBirth"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Gender      string `json:"gender"`
}

type BirthdayGreeting struct {
	Title   string `xml:"title"`
	Content string `xml:"content"`
}
