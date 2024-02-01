package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/repository"
)

type StudentRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *StudentRepositoryImpl) GetStudentByID(c context.Context, id uuid.UUID) (*models.Student, error) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE id = $1", id)

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&student)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("Errors in rows")
	}

	return &student, nil
}

func (r *StudentRepositoryImpl) GetStudentByEmail(c context.Context, email string) (*models.Student, error) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE email = $1", email)

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&student)
		if err != nil {
			return nil, err
		}
	}

	if student.Email == "" {
		return nil, errors.New("Email not found!")
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("Errors in rows")
	}

	return &student, nil
}

func (r *StudentRepositoryImpl) GetStudentByName(c context.Context, name string) ([]*models.Student, error) {
	var students []*models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM students WHERE name LIKE $1", name)

	defer rows.Close()

	for rows.Next() {
		student := new(models.Student)
		err := rows.StructScan(&student)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("Errors in rows")
	}

	return students, nil
}

func (r *StudentRepositoryImpl) CreateStudent(c context.Context, student *models.Student) error {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO students VALUES(:id, :name, :nim, :email, :password)", student)
	if err != nil {
		return errors.New("Failed to create student")
	}

	return nil
}

func (r *StudentRepositoryImpl) DeleteStudent(c context.Context, id uuid.UUID) error {
	_, err := r.DB.NamedExecContext(c, "DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return errors.New("Failed to delete student")
	}

	return nil
}

func NewStudentRepository(db *sqlx.DB) repository.StudentRepository {
	return &StudentRepositoryImpl{
		DB: db,
	}
}
