package users

import "time"

const (
	GenderMale   = 'M'
	GenderFemale = 'F'
	GenderOther  = 'O'
)

type User struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Email     string    `json:"email" dynamodbav:"email"`
	Password  string    `json:"password" dynamodbav:"password"`
	Profile   Profile   `json:"profile" dynamodbav:"profile"`
	LastLogin time.Time `json:"last_login" dynamodbav:"last_login"`
}

type Profile struct {
	Name      string    `json:"name" dynamodbav:"name"`
	Age       int       `json:"age" dynamodbav:"age"`
	BirthDate time.Time `json:"birth_date" dynamodbav:"birth_date"`
	Gender    rune      `json:"gender" dynamodbav:"gender"`
}

//Transport Models

type UserTransport struct {
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
	Password  string `json:"password"`
}

type CreateUserResponse struct {
	User UserTransport `json:"user"`
}

type GetUserResponse struct {
	User UserTransport `json:"user"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUserResponse struct {
	User UserTransport `json:"user"`
	JWT  JWT           `json:"jwt"`
}
