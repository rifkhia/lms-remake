package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/pkg"
	"github.com/rifkhia/lms-remake/internal/utils"
)

type ClassRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *ClassRepositoryImpl) CheckStudentClassExists(c context.Context, classId int, studentId uuid.UUID) (bool, pkg.CustomError) {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM student_class WHERE student_id = $1 AND class_id = $2 AND deleted_at IS NULL", studentId, classId)
	if err != nil {
		return false, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return count > 0, pkg.CustomError{}
}

func (r *ClassRepositoryImpl) GetClassByID(c context.Context, id int) (*models.Class, pkg.CustomError) {
	var class models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT id, name, description, key, teacher_id AS teacherid FROM classes WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, pkg.CustomError{
			Cause:   errors.New("no class with that id"),
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	err = rows.StructScan(&class)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return &class, pkg.CustomError{}
}

func (r *ClassRepositoryImpl) GetClassByTeacherID(c context.Context, teacherId int) ([]*models.Class, pkg.CustomError) {
	var classes []*models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM classes WHERE teacher_id = $1 AND deleted_at IS NULL", teacherId)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	defer rows.Close()

	for rows.Next() {
		class := new(models.Class)
		err := rows.StructScan(&class)
		if err != nil {
			return nil, pkg.CustomError{
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
				Code:    utils.INTERNAL_SERVER_ERROR,
			}
		}

		classes = append(classes, class)
	}

	err = rows.Err()

	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return classes, pkg.CustomError{}
}

func (r *ClassRepositoryImpl) GetClassByName(c context.Context, name string) ([]*models.Class, pkg.CustomError) {
	var classes []*models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT id, id, name, description, key, teacher_id AS teacherid FROM classes WHERE name LIKE '%'||$1||'%' AND deleted_at IS NULL", name)
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	defer rows.Close()

	for rows.Next() {
		class := new(models.Class)
		err = rows.StructScan(&class)
		if err != nil {
			return nil, pkg.CustomError{
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
				Code:    utils.INTERNAL_SERVER_ERROR,
			}
		}

		classes = append(classes, class)
	}

	err = rows.Err()

	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return classes, pkg.CustomError{}
}

func (r *ClassRepositoryImpl) CreateClass(c context.Context, class *models.Class) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO classes VALUES(DEFAULT, :name, :description, :key, :teacherid, now(), now(), null, :day, :starttime, :endtime)", class)
	if err != nil {
		return pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return pkg.CustomError{}
}

func (r *ClassRepositoryImpl) JoinClass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError {
	tempMap := map[string]interface{}{
		"classId":   classId,
		"studentId": studentId,
	}
	_, err := r.DB.NamedExecContext(c, "INSERT INTO student_class(student_id, class_id, created_at) VALUES(:studentId, :classId, now())", tempMap)
	if err != nil {
		return pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return pkg.CustomError{}
}

func (r *ClassRepositoryImpl) LeftCLass(c context.Context, classId int, studentId uuid.UUID) pkg.CustomError {
	tempMap := map[string]interface{}{
		"classId":   classId,
		"studentId": studentId,
	}
	_, err := r.DB.NamedExecContext(c, "UPDATE student_class SET deleted_at = now() WHERE student_id = :studentid AND class_id = :classid", tempMap)
	if err != nil {
		return pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return pkg.CustomError{}
}

func (r *ClassRepositoryImpl) DeleteClass(c context.Context, id int) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "UPDATE classes SET deleted_at = now() WHERE id = ?", id)
	if err != nil {
		return pkg.CustomError{
			Cause:   err,
			Service: utils.REPOSITORY_SERVICE,
			Code:    utils.INTERNAL_SERVER_ERROR,
		}
	}

	return pkg.CustomError{}
}

func (r *ClassRepositoryImpl) CreateClassSection(c context.Context, classSection *models.SectionClass) pkg.CustomError {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO class_sections VALUES (DEFAULT, :title, :description, :order, :classid, now(), now(), null)", classSection)
	if err != nil {
		customError := pkg.CustomError{
			Cause:   err,
			Code:    utils.INTERNAL_SERVER_ERROR,
			Service: utils.REPOSITORY_SERVICE,
		}
		return customError
	}

	return pkg.CustomError{}
}

func (r *ClassRepositoryImpl) GetClassSectionByClassId(c context.Context, classId int) ([]*models.SectionClass, pkg.CustomError) {
	var classSections []*models.SectionClass
	rows, err := r.DB.QueryxContext(c, "SELECT id, title, description, class_id AS classid, \"order\" FROM class_sections WHERE class_id = $1 AND deleted_at IS NULL", classId)
	if err != nil {
		customError := pkg.CustomError{
			Cause:   err,
			Code:    utils.INTERNAL_SERVER_ERROR,
			Service: utils.REPOSITORY_SERVICE,
		}
		return nil, customError
	}

	for rows.Next() {
		classSection := new(models.SectionClass)
		err = rows.StructScan(&classSection)
		if err != nil {
			return nil, pkg.CustomError{
				Cause:   err,
				Service: utils.REPOSITORY_SERVICE,
				Code:    utils.INTERNAL_SERVER_ERROR,
			}
		}

		classSections = append(classSections, classSection)
	}

	err = rows.Err()
	if err != nil {
		return nil, pkg.CustomError{
			Cause:   err,
			Code:    utils.INTERNAL_SERVER_ERROR,
			Service: utils.REPOSITORY_SERVICE,
		}
	}

	return classSections, pkg.CustomError{}
}

//func (r *ClassRepositoryImpl) InsertSubmissionTeacher(c context.Context, classSectionId int, linkFile string) pkg.CustomError {
//	_, err := r.DB.NamedExecContext("INSERT INTO submissions VALUES (DEFAULT, )")
//}

func NewClassRepository(db *sqlx.DB) ClassRepository {
	return &ClassRepositoryImpl{
		DB: db,
	}
}
