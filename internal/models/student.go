package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
	"net/mail"
	"time"
)

type Student struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NIM      int       `json:"NIM"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type StudentProfile struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	NIM         int       `json:"NIM"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
}

type StudentProfileRequest struct {
	ID          uuid.UUID `json:"ID"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
}

type StudentClass struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type StudentSchedule struct {
	ClassName string    `json:"class_name"`
	Day       string    `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type StudentLoginRequest struct {
	NIM      string `json:"NIM"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StudentRegisterRequest struct {
	Name     string `json:"name"`
	NIM      int    `json:"NIM"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *StudentProfileRequest) UpdateStudent(student *StudentProfileRequest) {
	if s.Address != "" {
		student.Address = s.Address
	}

	if s.DateOfBirth.IsZero() {
		student.DateOfBirth = s.DateOfBirth
	}

	if s.Gender != "" {
		student.Gender = s.Gender
	}

	if s.Phone != "" {
		student.Phone = s.Phone
	}

}

func (s *StudentRegisterRequest) NewStudent() (*Student, pkg.CustomError) {
	studentId := uuid.New()

	//validate name
	if s.Name == "" {
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("name cant be blank"),
		}
	}

	//validate email
	_, err := mail.ParseAddress(s.Email)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("invalidate email"),
		}
	}

	//validate npm
	if s.NIM == 0 {
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("nim cant be blank"),
		}
	}

	//validate password
	if len(s.Password) < 8 {
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Service: "models",
			Cause:   errors.New("password length must be more than 8"),
		}
	}

	//hashing password
	s.Password, err = utils.GeneratePassword(s.Password)

	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Service: "models",
			Cause:   err,
		}
	}

	return &Student{
		ID:       studentId,
		Name:     s.Name,
		NIM:      s.NIM,
		Email:    s.Email,
		Password: s.Password,
	}, pkg.CustomError{}
}
