package usecase

import (
	"context"
	"github.com/rifkhia/lms-remake/internal/models"
	error2 "github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
	"log"
)

type teacherUsecaseImpl struct {
	teacherRepo repository.TeacherRepository
}

type TeacherUsecase interface {
	LoginTeacher(c context.Context, request *models.TeacherLoginRequest) (interface{}, error2.CustomError)
	RegisterTeacher(c context.Context, request *models.TeacherRegisterRequest) (interface{}, error2.CustomError)
}

func (s *teacherUsecaseImpl) LoginTeacher(c context.Context, request *models.TeacherLoginRequest) (interface{}, error2.CustomError) {
	teacherResult, err := s.teacherRepo.GetTeacherByEmail(c, request.Email)
	if err.Cause != nil {
		return nil, err
	}

	log.Println(request.Email)

	err = utils.ValidatePassword(teacherResult.Password, request.Password)
	if err.Cause != nil {
		return nil, err
	}

	accessToken, err := utils.CreateAccessToken(teacherResult.ID, utils.TEACHER_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	refeshToken, err := utils.CreateRefreshToken(teacherResult.ID, accessToken, utils.TEACHER_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access token":  accessToken,
		"refresh token": refeshToken,
	}, error2.CustomError{}
}

func (s *teacherUsecaseImpl) RegisterTeacher(c context.Context, request *models.TeacherRegisterRequest) (interface{}, error2.CustomError) {
	teacherRequest, err := request.NewTeacher()
	if err.Cause != nil {
		return nil, err
	}

	err = s.teacherRepo.CreateTeacher(c, teacherRequest)
	if err.Cause != nil {
		return nil, err
	}

	accessToken, err := utils.CreateAccessToken(teacherRequest.ID, utils.TEACHER_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	refeshToken, err := utils.CreateRefreshToken(teacherRequest.ID, accessToken, utils.TEACHER_ROLE)
	if err.Cause != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access token":  accessToken,
		"refresh token": refeshToken,
	}, err
}

func NewTeacherUsecase(repo repository.TeacherRepository) TeacherUsecase {
	return &teacherUsecaseImpl{
		teacherRepo: repo,
	}
}
