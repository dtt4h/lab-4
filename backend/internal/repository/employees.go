package repository

import (
	"context"

	"lab-4/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeesRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	GetByID(ctx context.Context, id string) (*models.Employee, error)
	List(ctx context.Context) ([]models.Employee, error)
	Update(ctx context.Context, employee *models.Employee) error
	Delete(ctx context.Context, id string) error
}

type employeesRepository struct{ db *pgxpool.Pool }

func NewEmployeesRepository(db *pgxpool.Pool) EmployeesRepository {
	return &employeesRepository{db: db}
}

func (r *employeesRepository) Create(ctx context.Context, employee *models.Employee) error {
	err := r.db.QueryRow(ctx, `INSERT INTO employees (user_id, hotel_id, full_name, position) VALUES ($1, $2, $3, $4) RETURNING id`, employee.UserID, employee.HotelID, employee.FullName, employee.Position).Scan(&employee.ID)
	return normalizeError(err)
}

func (r *employeesRepository) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	employee := new(models.Employee)
	err := r.db.QueryRow(ctx, `SELECT id, user_id, hotel_id, full_name, position FROM employees WHERE id = $1`, id).Scan(&employee.ID, &employee.UserID, &employee.HotelID, &employee.FullName, &employee.Position)
	if err != nil {
		return nil, normalizeError(err)
	}
	return employee, nil
}

func (r *employeesRepository) List(ctx context.Context) ([]models.Employee, error) {
	rows, err := r.db.Query(ctx, `SELECT id, user_id, hotel_id, full_name, position FROM employees ORDER BY full_name`)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	employees := make([]models.Employee, 0)
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.UserID, &employee.HotelID, &employee.FullName, &employee.Position); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeesRepository) Update(ctx context.Context, employee *models.Employee) error {
	result, err := r.db.Exec(ctx, `UPDATE employees SET user_id = $2, hotel_id = $3, full_name = $4, position = $5 WHERE id = $1`, employee.ID, employee.UserID, employee.HotelID, employee.FullName, employee.Position)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *employeesRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return normalizeError(err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
