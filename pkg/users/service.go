package users

import (
	"context"
	"net/mail"
	"time"

	"github.com/eduardohoraciosanto/users/internal/db"
	"github.com/eduardohoraciosanto/users/internal/errors"
	"github.com/eduardohoraciosanto/users/internal/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, name string, age int, email string, birthdate time.Time, gender rune, pass string) (User, error)
	Login(ctx context.Context, email, password string) (User, error)
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

func (u *users) Create(ctx context.Context, name string, age int, email string, birthdate time.Time, gender rune, pass string) (User, error) {
	userDB := &User{}
	//check if user exists
	err := u.db.Get(ctx, email, userDB)
	if err == nil {
		u.log.WithField("email", email).Error(ctx, "email already taken")
		return User{}, errors.New(errors.EmailTakenCode, "email already used")
	}

	switch gender {
	case GenderFemale, GenderMale, GenderOther:
		break
	default:
		u.log.WithField("gender", gender).Error(ctx, "gender not supported")
		return User{}, errors.New(errors.BadGenderCode, "gender not supported")
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		u.log.WithError(err).
			WithField("email", email).Error(ctx, "email is not valid")
		return User{}, errors.NewFromError(errors.BadEmailCode, err)
	}

	id := uuid.NewString()

	//Hash password
	hPass, err := bcrypt.GenerateFromPassword([]byte(id+pass), 10)
	if err != nil {
		u.log.WithError(err).
			WithField("email", email).Error(ctx, "unable to hash password")
		return User{}, errors.NewFromError(errors.InternalErrorCode, err)
	}

	user := User{
		ID:       id,
		Email:    email,
		Password: string(hPass),
		Profile: Profile{
			Name:      name,
			Age:       age,
			BirthDate: birthdate,
			Gender:    gender,
		},
	}

	err = u.Save(ctx, user)
	if err != nil {
		u.log.WithError(err).
			WithField("email", email).Error(ctx, "unable to save new user")
		return User{}, errors.NewFromError(errors.InternalErrorCode, err)
	}

	u.log.WithField("user_id", id).Info(ctx, "User created")

	return user, nil

}

func (u *users) Login(ctx context.Context, email, password string) (User, error) {
	user, err := u.Get(ctx, email)
	if err != nil {
		u.log.WithField("user_email", email).
			WithError(err).Error(ctx, "Unable to get user")
		return User{}, errors.New(errors.InvalidCredentialsCode, "Invalid Credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.ID+password))
	if err != nil {
		u.log.WithField("user_email", email).
			WithError(err).Error(ctx, "Invalid Password")
		return User{}, errors.New(errors.InvalidCredentialsCode, "Invalid Credentials")
	}

	user.LastLogin = time.Now()

	err = u.Save(ctx, user)
	if err != nil {
		u.log.WithField("user_email", email).
			WithError(err).Error(ctx, "error saving last login")
		return User{}, errors.NewFromError(errors.InternalErrorCode, err)
	}

	return user, nil
}

func (u *users) Get(ctx context.Context, email string) (User, error) {

	user := User{}

	err := u.db.Get(ctx, email, &user)
	if err != nil {
		u.log.WithField("user_email", email).
			WithError(err).Error(ctx, "Unable to get user")
		return User{}, errors.NewFromError(errors.UserNotFoundCode, err)
	}

	u.log.WithField("user_email", email).Info(ctx, "Got User from DB")

	return user, nil

}
func (u *users) Save(ctx context.Context, user User) error {

	err := u.db.Set(ctx, user.Email, user)
	if err != nil {
		u.log.WithError(err).
			WithField("user_email", user.Email).
			Error(ctx, "Unable to save user on DB")
		return errors.NewFromError(errors.DBErrorSavingCode, err)
	}

	u.log.WithField("user_email", user.Email).Info(ctx, "User saved successfully")

	return nil
}

func (u *users) Delete(ctx context.Context, user User) error {
	err := u.db.Delete(ctx, user.Email)
	if err != nil {
		u.log.WithError(err).
			WithField("user_email", user.Email).
			Error(ctx, "Unable to delete user from DB")
		return errors.NewFromError(errors.DBErrorDeletingCode, err)
	}

	u.log.WithField("user_email", user.Email).Info(ctx, "User deleted successfully")

	return nil
}
