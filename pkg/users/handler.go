package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/response"
)

type Handler struct {
	Service Service
}

//Create is the handler for creating a user
func (c *Handler) Create(w http.ResponseWriter, r *http.Request) {
	//parse body into our model
	payload := &CreateUserRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(payload)
	if err != nil {
		response.RespondWithError(w, errors.NewFromError(errors.ParsingErrorCode, err))
		return
	}

	//parse birthdate
	bd, err := time.Parse("2006-01-02", payload.Birthdate)
	if err != nil {
		response.RespondWithError(w, errors.NewFromError(errors.ParsingErrorCode, err))
		return
	}

	//parse gender
	if len(payload.Gender) >= 2 {
		response.RespondWithError(w, errors.New(errors.ParsingErrorCode, "Incorrect Gender. Should be one of M, F or O"))
		return
	}
	gender := rune(payload.Gender[0])

	user, err := c.Service.Create(r.Context(), payload.Name, payload.Age, payload.Email, bd, gender)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	res := CreateUserResponse{
		User: UserTransport{
			ID:        user.ID,
			Name:      user.Name,
			Age:       user.Age,
			Email:     user.Email,
			Birthdate: user.Birthdate.Format("2006-01-02"),
			Gender:    string(user.Gender),
		},
	}
	response.RespondWithData(w, http.StatusOK, res)
}
