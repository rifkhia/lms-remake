package teacher_impl

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/rifkhia/lms-remake/internal/models"
)

type TeacherRepositoryImpl struct {
	DB *sqlx.DB
}

func (r *TeacherRepositoryImpl) GetTeacherByEmail(c context.Context, email string) (*models.Teacher, error) {
	var teacher models.Teacher

	rows, err := r.DB.QueryxContext(c, "SELECT * FROM teachers WHERE email = $1", email)

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&teacher)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()

	if err != nil {
		return nil, errors.New("Errors in rows")
	}

	return &teacher, nil
}
