package service

import (
	"context"
	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type EmployeesService interface {
	Create(context.Context, *models.Employee) error
	GetByID(context.Context, string) (*models.Employee, error)
	List(context.Context) ([]models.Employee, error)
	Update(context.Context, *models.Employee) error
	Delete(context.Context, string) error
}
type employeesService struct {
	employees repository.EmployeesRepository
}

func NewEmployeesService(employees repository.EmployeesRepository) EmployeesService {
	return &employeesService{employees: employees}
}
func (s *employeesService) Create(ctx context.Context, employee *models.Employee) error {
	if err := required(employee.HotelID, employee.FullName, employee.Position); err != nil {
		return err
	}
	return s.employees.Create(ctx, employee)
}
func (s *employeesService) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.employees.GetByID(ctx, id)
}
func (s *employeesService) List(ctx context.Context) ([]models.Employee, error) {
	return s.employees.List(ctx)
}
func (s *employeesService) Update(ctx context.Context, employee *models.Employee) error {
	if err := required(employee.ID, employee.HotelID, employee.FullName, employee.Position); err != nil {
		return err
	}
	return s.employees.Update(ctx, employee)
}
func (s *employeesService) Delete(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.employees.Delete(ctx, id)
}
