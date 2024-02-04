package models

import (
	"errors"
	"github.com/google/uuid"
	error2 "github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
	"net/mail"
)

type Teacher struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NPM      int       `json:"NPM"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type TeacherLoginRequest struct {
	NPM      int    `json:"NPM"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TeacherRegisterRequest struct {
	Name     string `json:"name"`
	NPM      int    `json:"NPM"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *TeacherRegisterRequest) NewTeacher() (*Teacher, error2.CustomError) {
	var custErr error2.CustomError
	//generating uuid
	studentId := uuid.New()

	//validate name
	if s.Name == "" {
		custErr = error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("name cant be blank"),
		}
		return nil, custErr
	}

	//validate email
	_, err := mail.ParseAddress(s.Email)
	if err != nil {
		custErr = error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("invalidate email"),
		}
		return nil, custErr
	}

	//validate npm
	if s.NPM == 0 {
		custErr = error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("npm cant be blank"),
		}
		return nil, custErr
	}

	//validate password
	if len(s.Password) < 8 {
		custErr = error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("password length must be more than 8"),
		}
		return nil, custErr
	}

	//hashing password
	s.Password, err = utils.GeneratePassword(s.Password)

	if err != nil {
		custErr = error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   err,
		}
		return nil, custErr
	}

	return &Teacher{
		ID:       studentId,
		Name:     s.Name,
		NPM:      s.NPM,
		Email:    s.Email,
		Password: s.Password,
	}, custErr
}
