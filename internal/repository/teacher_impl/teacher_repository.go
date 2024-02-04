package teacher_impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
	error2 "github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
	"strings"
)

type TeacherRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *TeacherRepositoryImpl) GetTeacherById(c context.Context, id uuid.UUID) (*models.Teacher, error2.CustomError) {
	var teacher models.Teacher

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM teachers WHERE id = $1 AND deleted_at IS NULL", id)

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, error2.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return nil, error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no teacher found using current id"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		err := rows.StructScan(&teacher)
		if err != nil {
			return nil, error2.CustomError{
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
				Code:    utils.INTERNAL_SERVER_ERROR,
			}
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, error2.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return &teacher, error2.CustomError{}
}

func (r *TeacherRepositoryImpl) GetTeacherByEmail(c context.Context, email string) (*models.Teacher, error2.CustomError) {
	var teacher models.Teacher

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM teachers WHERE email = $1 AND deleted_at IS NULL", email)

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, error2.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return nil, error2.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no teacher found using current email"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		err := rows.StructScan(&teacher)
		if err != nil {
			return nil, error2.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, error2.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return &teacher, error2.CustomError{}
}

func (r *TeacherRepositoryImpl) CreateTeacher(c context.Context, request *models.Teacher) error2.CustomError {
	var custErr error2.CustomError
	_, err := r.DB.NamedExecContext(c, "INSERT INTO teachers VALUES (:id, :name, :npm, :email, :password, now(), now(), null)", request)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return error2.CustomError{
				Code:    utils.BAD_REQUEST,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		custErr = error2.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Service: utils.REPOSITORY_SERVICE,
			Cause:   err,
		}
		return custErr
	}

	return error2.CustomError{}
}

func NewTeacherRepository(db *sqlx.DB) repository.TeacherRepository {
	return &TeacherRepositoryImpl{
		DB: db,
	}
}
