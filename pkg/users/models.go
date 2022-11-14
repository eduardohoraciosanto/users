package users

import "time"

const (
	GenderMale   = 'M'
	GenderFemale = 'F'
	GenderOther  = 'O'
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Email     string    `json:"email"`
	Birthdate time.Time `json:"birthdate"`
	Gender    rune      `json:"gender"`
}

//Transport Models

type UserTransport struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"`
}

type CreateUserResponse struct {
	User UserTransport `json:"user"`
}
