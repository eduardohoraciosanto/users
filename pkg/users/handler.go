package users

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/response"
	"github.com/gorilla/mux"
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

	user, err := c.Service.Create(r.Context(), payload.Name, payload.Age, payload.Email, bd, gender, payload.Password)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	res := CreateUserResponse{
		User: UserTransport{
			ID:        user.ID,
			Name:      user.Profile.Name,
			Age:       user.Profile.Age,
			Email:     user.Profile.Email,
			Birthdate: user.Profile.BirthDate.Format("2006-01-02"),
			Gender:    string(user.Profile.Gender),
		},
	}
	response.RespondWithData(w, http.StatusOK, res)
}

//Get retrieves a user from the DB if it exists
func (c *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, err := c.Service.Get(r.Context(), vars["user_email"])
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	res := GetUserResponse{
		User: UserTransport{
			ID:        user.ID,
			Name:      user.Profile.Name,
			Age:       user.Profile.Age,
			Email:     user.Profile.Email,
			Birthdate: user.Profile.BirthDate.Format("2006-01-02"),
			Gender:    string(user.Profile.Gender),
		},
	}
	response.RespondWithData(w, http.StatusOK, res)
}

func (c *Handler) Login(w http.ResponseWriter, r *http.Request) {
	payload := &LoginUserRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(payload)
	if err != nil {
		response.RespondWithError(w, errors.NewFromError(errors.ParsingErrorCode, err))
		return
	}

	user, err := c.Service.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	res := LoginUserResponse{
		User: UserTransport{
			ID:        user.ID,
			Name:      user.Profile.Name,
			Age:       user.Profile.Age,
			Email:     user.Profile.Email,
			Birthdate: user.Profile.BirthDate.Format("2006-01-02"),
			Gender:    string(user.Profile.Gender),
		},
		JWT: JWT{
			AccessToken:  "TBD",
			RefreshToken: "TBD",
		},
	}

	response.RespondWithData(w, http.StatusOK, res)
}
