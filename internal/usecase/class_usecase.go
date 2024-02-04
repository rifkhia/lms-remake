package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type ClassUsecase interface {
	FetchClassById(c context.Context, id int) (*models.Class, pkg.CustomError)
	FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError)
	FetchClassByName(c context.Context, name string) ([]*models.Class, pkg.CustomError)
	CreateClass(c context.Context, request *models.Class) pkg.CustomError
	JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) pkg.CustomError
	CreateSectionClass(c context.Context, request *models.SectionClass) pkg.CustomError
}

type classUsecaseImpl struct {
	classRepo repository.ClassRepository
}

func (s *classUsecaseImpl) FetchClassById(c context.Context, id int) (*models.Class, pkg.CustomError) {
	classResult, err := s.classRepo.GetClassByID(c, id)
	if err.Cause != nil {
		return nil, err
	}

	classSection, err := s.classRepo.GetClassSectionByClassId(c, id)
	if err.Cause != nil {
		return nil, err
	}

	classResult.ClassSection = classSection

	return classResult, pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError) {
	classResult, err := s.classRepo.GetClassByTeacherID(c, teacherId)
	if err.Cause != nil {
		return nil, err
	}

	return classResult, pkg.CustomError{}
}

func (s *classUsecaseImpl) FetchClassByName(c context.Context, name string) ([]*models.Class, pkg.CustomError) {
	classResult, err := s.classRepo.GetClassByName(c, name)
	if err.Cause != nil {
		return nil, err
	}

	return classResult, pkg.CustomError{}
}

func (s *classUsecaseImpl) CreateClass(c context.Context, request *models.Class) pkg.CustomError {
	request.Key = utils.GenerateClassKey(8)

	err := s.classRepo.CreateClass(c, request)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) pkg.CustomError {
	exists, err := s.classRepo.CheckStudentClassExists(c, classId, studentId)
	if err.Cause != nil {
		return err
	}

	if exists {
		return pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("you already join this class"),
			Service: utils.USECASE_SERVICE,
		}
	}

	classResult, err := s.classRepo.GetClassByID(c, classId)
	if err.Cause != nil {
		return err
	}

	if classResult.Key != key {
		return pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("missmatch class key"),
			Service: utils.USECASE_SERVICE,
		}
	}

	err = s.classRepo.JoinClass(c, classId, studentId)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func (s *classUsecaseImpl) CreateSectionClass(c context.Context, request *models.SectionClass) pkg.CustomError {
	err := s.classRepo.CreateClassSection(c, request)
	if err.Cause != nil {
		return err
	}

	return pkg.CustomError{}
}

func NewClassUsecase(repo repository.ClassRepository) ClassUsecase {
	return &classUsecaseImpl{
		classRepo: repo,
	}
}
