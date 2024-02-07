package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/dto"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
	"strings"
)

type StudentRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *StudentRepositoryImpl) GetStudentByID(c context.Context, id uuid.UUID) (*models.Student, pkg.CustomError) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT id, name, nim, email, password FROM students WHERE id = $1 AND deleted_at IS NULL", id)
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

	err = rows.StructScan(&student)
	if err != nil {
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

func (r *StudentRepositoryImpl) GetStudentByEmail(c context.Context, email string) (*models.Student, pkg.CustomError) {
	var student models.Student

	rows, err := r.DB.QueryxContext(c, "SELECT id, name, nim, email, password FROM students WHERE email = $1 AND deleted_at IS NULL", email)
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

	err = rows.StructScan(&student)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
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

	return students, pkg.CustomError{}
}

func (r *StudentRepositoryImpl) GetStudentByClassId(c context.Context, classId int) ([]*models.StudentClass, pkg.CustomError) {
	var students []*models.StudentClass
	rows, err := r.DB.QueryxContext(c, "SELECT students.id, students.name FROM students INNER JOIN public.student_class sc on students.id = sc.student_id WHERE sc.class_id = $1", classId)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	for rows.Next() {
		student := new(models.StudentClass)
		err = rows.StructScan(&student)
		if err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		students = append(students, student)
	}

	return students, pkg.CustomError{}
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

	return pkg.CustomError{}
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

	return pkg.CustomError{}
}

func (r *StudentRepositoryImpl) GetStudentProfile(c context.Context, id uuid.UUID) (*dto.StudentProfileRequest, pkg.CustomError) {
	var student dto.StudentProfileRequest
	rows, err := r.DB.QueryxContext(c, "SELECT id, dateofbirth, gender, address, phone FROM student_profile WHERE id = $1", id)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	if !rows.Next() {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   errors.New("no student profile found"),
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	err = rows.StructScan(&student)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return &student, pkg.CustomError{}
}

func (r *StudentRepositoryImpl) AddStudentProfile(c context.Context, student *dto.StudentProfileRequest) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO student_profile VALUES(:id, :dateofbirth, :gender, :address, :phone, now(), now(), null)", student)
	if err != nil {
		return pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return pkg.CustomError{}
}

func (r *StudentRepositoryImpl) EditStudentProfile(c context.Context, student *dto.StudentProfileRequest) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "UPDATE student_profile SET dateofbirth=:dateofbirth, gender=:gender, address=:address, phone=:phone, updated_at=now() WHERE id=:id", student)
	if err != nil {
		return pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return pkg.CustomError{}
}

func (r *StudentRepositoryImpl) FetchStudentClass(c context.Context, id uuid.UUID) ([]*models.StudentSchedule, pkg.CustomError) {
	var schedules []*models.StudentSchedule
	rows, err := r.DB.QueryxContext(c, "SELECT c.name, c.day, c.start_time as starttime, c.end_time as endtime FROM students LEFT JOIN public.student_class sc on students.id = sc.student_id LEFT JOIN public.classes c on c.id = sc.class_id WHERE students.id=$1 ORDER BY day, start_time", id)
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	defer rows.Close()

	for rows.Next() {
		var schedule = new(models.StudentSchedule)
		err = rows.StructScan(&schedule)
		if err != nil {
			return nil, pkg.CustomError{
				Code:    utils.INTERNAL_SERVER_ERROR,
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
			}
		}
		schedules = append(schedules, schedule)
	}

	err = rows.Err()
	if err != nil {
		return nil, pkg.CustomError{
			Code:    utils.INTERNAL_SERVER_ERROR,
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return schedules, pkg.CustomError{}
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	return &StudentRepositoryImpl{
		DB: db,
	}
}
