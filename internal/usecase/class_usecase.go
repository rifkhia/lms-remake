package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type ClassUsecase interface {
	FetchClassById(c context.Context, id int) (*models.Class, error)
	FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, error)
	FetchClassByName(c context.Context, name string) ([]*models.Class, error)
	CreateClass(c context.Context, request *models.Class) error
	JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) error
}

type classUsecaseImpl struct {
	classRepo repository.ClassRepository
}

func (s *classUsecaseImpl) FetchClassById(c context.Context, id int) (*models.Class, error) {
	classResult, err := s.classRepo.GetClassByID(c, id)
	if err != nil {
		return nil, err
	}

	return classResult, nil
}

func (s *classUsecaseImpl) FetchClassByTeacherId(c context.Context, teacherId int) ([]*models.Class, error) {
	classResult, err := s.classRepo.GetClassByTeacherID(c, teacherId)
	if err != nil {
		return nil, err
	}

	return classResult, nil
}

func (s *classUsecaseImpl) FetchClassByName(c context.Context, name string) ([]*models.Class, error) {
	classResult, err := s.classRepo.GetClassByName(c, name)
	if err != nil {
		return nil, err
	}

	return classResult, nil
}

func (s *classUsecaseImpl) CreateClass(c context.Context, request *models.Class) error {
	request.Key = utils.GenerateClassKey(8)

	err := s.classRepo.CreateClass(c, request)
	if err != nil {
		return err
	}

	return nil
}

func (s *classUsecaseImpl) JoinClass(c context.Context, studentId uuid.UUID, classId int, key string) error {
	exists, err := s.classRepo.CheckStudentClassExists(c, classId, studentId)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("you already join this class!")
	}

	classResult, err := s.classRepo.GetClassByID(c, classId)
	if err != nil {
		return err
	}

	if classResult.Key != key {
		return errors.New("Wrong class key")
	}

	err = s.classRepo.JoinClass(c, classId, studentId)
	if err != nil {
		return err
	}

	return nil
}

func NewClassUsecase(repo repository.ClassRepository) ClassUsecase {
	return &classUsecaseImpl{
		classRepo: repo,
	}
}
