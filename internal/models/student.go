package models

import (
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type Student struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NIM      string    `json:"NIM"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

type LoginRegisterRequest struct {
	Name     string `json:"name"`
	NIM      string `json:"NIM"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *LoginRegisterRequest) NewStudent() (*Student, error) {
	var err error
	studentId := uuid.New()
	s.Password, err = utils.GeneratePassword(s.Password)

	if err != nil {
		return nil, err
	}

	return &Student{
		ID:       studentId,
		Name:     s.Name,
		NIM:      s.NIM,
		Email:    s.Email,
		Password: s.Password,
	}, nil
}
