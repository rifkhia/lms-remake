package class_impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
	"github.com/rifkhia/lms-remake/internal/repository"
)

type ClassRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *ClassRepositoryImpl) CheckStudentClassExists(c context.Context, classId int, studentId uuid.UUID) (bool, error) {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM student_class WHERE student_id = $1 AND class_id = $2", studentId, classId)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ClassRepositoryImpl) GetClassByID(c context.Context, id int) (*models.Class, error) {
	var class models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT id, name, description, key FROM classes WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&class)
		if err != nil {
			panic(err)

			return nil, errors.New("error in structscan")
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("error in row")
	}

	return &class, nil
}

func (r *ClassRepositoryImpl) GetClassByTeacherID(c context.Context, teacherId int) ([]*models.Class, error) {
	var classes []*models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM classes WHERE teacher_id = $1", teacherId)

	defer rows.Close()

	for rows.Next() {
		class := new(models.Class)
		err := rows.StructScan(&class)
		if err != nil {
			return nil, errors.New("error in structscan")
		}

		classes = append(classes, class)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("error in row")
	}

	return classes, nil
}

func (r *ClassRepositoryImpl) GetClassByName(c context.Context, name string) ([]*models.Class, error) {
	var classes []*models.Class

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM classes WHERE name LIKE $1", name)

	defer rows.Close()

	for rows.Next() {
		class := new(models.Class)
		err := rows.StructScan(&class)
		if err != nil {
			return nil, errors.New("error in structscan")
		}

		classes = append(classes, class)
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("error in row")
	}

	return classes, nil
}

func (r *ClassRepositoryImpl) CreateClass(c context.Context, class *models.Class) error {
	_, err := r.DB.NamedExecContext(c, "INSERT INTO classes VALUES(:id, :name, :description, :key, :teacher_id)", class)
	if err != nil {
		return errors.New("Failed to create class")
	}

	return nil
}

func (r *ClassRepositoryImpl) JoinClass(c context.Context, classId int, studentId uuid.UUID) error {
	tempMap := map[string]interface{}{
		"classId":   classId,
		"studentId": studentId,
	}
	_, err := r.DB.NamedExecContext(c, "INSERT INTO student_class(student_id, class_id) VALUES(:studentId, :classId)", tempMap)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClassRepositoryImpl) DeleteClass(c context.Context, id int) error {
	_, err := r.DB.NamedExecContext(c, "DELETE FROM classes WHERE id = ?", id)
	if err != nil {
		return errors.New("Failed to delete class")
	}

	return nil
}

func NewClassRepository(db *sqlx.DB) repository.ClassRepository {
	return &ClassRepositoryImpl{
		DB: db,
	}
}
