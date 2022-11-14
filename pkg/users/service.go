package users

import (
	"context"
	"encoding/json"
	"net/mail"
	"time"

	"github.com/eduardohoraciosanto/users/internal/db"
	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/logger"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, name string, age int, email string, birthdate time.Time, gender rune) (User, error)
	Get(ctx context.Context, id string) (User, error)
	Save(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
}

type users struct {
	db  db.DB
	log logger.Logger
}

func NewService(db db.DB, log logger.Logger) Service {
	return &users{
		db:  db,
		log: log,
	}
}

func (u *users) Create(ctx context.Context, name string, age int, email string, birthdate time.Time, gender rune) (User, error) {

	switch gender {
	case GenderFemale, GenderMale, GenderOther:
		break
	default:
		u.log.WithField("gender", gender).Error(ctx, "gender not supported")
		return User{}, errors.New(errors.BadGenderCode, "gender not supported")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		u.log.WithError(err).
			WithField("email", email).Error(ctx, "email is not valid")
		return User{}, errors.NewFromError(errors.BadEmailCode, err)
	}

	id := uuid.NewString()

	u.log.WithField("user_id", id).Info(ctx, "User created")

	return User{
		ID:        id,
		Name:      name,
		Age:       age,
		Email:     email,
		Birthdate: birthdate,
		Gender:    gender,
	}, nil

}
func (u *users) Get(ctx context.Context, id string) (User, error) {

	userDB, err := u.db.Get(ctx, id)
	if err != nil {
		u.log.WithError(err).Error(ctx, "Unable to get user")
		return User{}, errors.NewFromError(errors.UserNotFoundCode, err)
	}

	userStr, ok := userDB.(string)
	if !ok {
		u.log.WithField("db_data", userDB).Error(ctx, "Unexpected data on DB")
		return User{}, errors.New(errors.InternalErrorCode, "unexpected data on DB")
	}

	user := User{}
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		u.log.WithError(err).
			WithField("raw_user", userStr).
			Error(ctx, "Unable parse saved user")
		return User{}, errors.NewFromError(errors.InternalErrorCode, err)
	}

	u.log.WithField("user_id", id).Info(ctx, "Got User from DB")

	return user, nil

}
func (u *users) Save(ctx context.Context, user User) error {

	userBytes, err := json.Marshal(user)
	if err != nil {
		u.log.WithError(err).
			WithField("user_id", user.ID).
			Error(ctx, "Unable to marshal user")
		return errors.NewFromError(errors.InternalErrorCode, err)
	}

	err = u.db.Set(ctx, user.ID, userBytes)
	if err != nil {
		u.log.WithError(err).
			WithField("user_id", user.ID).
			Error(ctx, "Unable to save user on DB")
		return errors.NewFromError(errors.DBErrorSavingCode, err)
	}

	u.log.WithField("user_id", user.ID).Info(ctx, "User saved successfully")

	return nil
}

func (u *users) Delete(ctx context.Context, user User) error {
	err := u.db.Delete(ctx, user.ID)
	if err != nil {
		u.log.WithError(err).
			WithField("user_id", user.ID).
			Error(ctx, "Unable to delete user from DB")
		return errors.NewFromError(errors.DBErrorDeletingCode, err)
	}

	u.log.WithField("user_id", user.ID).Info(ctx, "User deleted successfully")

	return nil
}
