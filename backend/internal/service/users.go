package service

import (
	"context"
	"strings"

	"lab-4/backend/internal/models"
	"lab-4/backend/internal/repository"
)

type UsersService interface {
	Register(ctx context.Context, params models.CreateUserParams) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id, role string) error
	List(ctx context.Context) ([]models.User, error)
}

type usersService struct{ users repository.UsersRepository }

func NewUsersService(users repository.UsersRepository) UsersService {
	return &usersService{users: users}
}
func (s *usersService) Register(ctx context.Context, params models.CreateUserParams) (*models.User, error) {
	if len(params.Username) < 3 || len(params.Username) > 50 || !strings.Contains(params.Email, "@") || params.PasswordHash == "" {
		return nil, ErrValidation
	}
	user := &models.User{Username: params.Username, Email: params.Email, PasswordHash: params.PasswordHash, Role: "guest"}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (s *usersService) GetByID(ctx context.Context, id string) (*models.User, error) {
	if err := required(id); err != nil {
		return nil, err
	}
	return s.users.GetByID(ctx, id)
}
func (s *usersService) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	if err := required(login); err != nil {
		return nil, err
	}
	return s.users.GetUserByLogin(ctx, login)
}
func (s *usersService) UpdateLastLogin(ctx context.Context, id string) error {
	if err := required(id); err != nil {
		return err
	}
	return s.users.UpdateLastLogin(ctx, id)
}
func (s *usersService) UpdateRole(ctx context.Context, id, role string) error {
	if err := required(id, role); err != nil {
		return err
	}
	switch role {
	case "super_admin", "booking_manager", "security_manager", "hotel_manager", "room_manager", "guest":
		return s.users.UpdateRole(ctx, id, role)
	}
	return ErrValidation
}
func (s *usersService) List(ctx context.Context) ([]models.User, error) { return s.users.List(ctx) }
