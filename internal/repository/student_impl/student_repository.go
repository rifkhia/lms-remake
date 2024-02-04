package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/repository"
	"github.com/rifkhia/lms-remake/internal/utils"
	"strings"
)

type StudentRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *StudentRepositoryImpl) GetStudentByID(c context.Context, id uuid.UUID) (*models.Student, pkg.CustomError) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no student found using current email"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		err := rows.StructScan(&student)
		if err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return &student, pkg.CustomError{}
}

func (r *StudentRepositoryImpl) GetStudentByEmail(c context.Context, email string) (*models.Student, pkg.CustomError) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE email = $1 AND deleted_at IS NULL", email)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no student found using current email"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		err := rows.StructScan(&student)
		if err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
	}

	if student.Email == "" {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return &student, pkg.CustomError{}
}

func (r *StudentRepositoryImpl) GetStudentByName(c context.Context, name string) ([]*models.Student, pkg.CustomError) {
	var students []*models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE name LIKE $1 AND deleted_at IS NULL", name)

	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return nil, pkg.CustomError{
			Code:    utils.BAD_REQUEST,
			Cause:   errors.New("no student found using current name"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		student := new(models.Student)
		err := rows.StructScan(&student)
		if err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}

		students = append(students, student)
	}

	err = rows.Err()

	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return students, pkg.CustomError{
		Code:    utils.INTERNAL_SERVER_ERROR,
		Cause:   err,
		Service: utils.REPOSITORY_SERVICE,
	}
}

func (r *StudentRepositoryImpl) CreateStudent(c context.Context, student *models.Student) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO students VALUES(:id, :name, :nim, :email, :password, now(), now(), null)", student)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return pkg.CustomError{
				Code:    utils.BAD_REQUEST,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		return pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return pkg.CustomError{
		Code:    utils.INTERNAL_SERVER_ERROR,
		Cause:   err,
		Service: utils.REPOSITORY_SERVICE,
	}
}

func (r *StudentRepositoryImpl) DeleteStudent(c context.Context, id uuid.UUID) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "UPDATE students SET deleted_at = now() WHERE id = ?", id)
	if err != nil {
		return pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return pkg.CustomError{
		Code:    utils.INTERNAL_SERVER_ERROR,
		Cause:   err,
		Service: utils.REPOSITORY_SERVICE,
	}
}

func NewStudentRepository(db *sqlx.DB) repository.StudentRepository {
	return &StudentRepositoryImpl{
		DB: db,
	}
}
