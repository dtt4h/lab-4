package repository

import (
	"context"
	"lab-4/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id, role string) error
	List(ctx context.Context) ([]models.User, error)
}

type usersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(db *pgxpool.Pool) UsersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) Create(ctx context.Context, user *models.User) error {
	const query = `INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, last_login_at, role`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.Role).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.Role,
	)
	return normalizeError(err)
}

func (r *usersRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	const query = `SELECT id, username, email, password_hash, created_at, updated_at, last_login_at, role
		FROM users WHERE id = $1`
	user := new(models.User)
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.Role)
	if err != nil {
		return nil, normalizeError(err)
	}
	return user, nil
}

func (r *usersRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	const query = `SELECT id, username, email, password_hash, created_at, updated_at, last_login_at, role
		FROM users WHERE username = $1 OR email = $1`
	user := new(models.User)
	err := r.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.Role)
	if err != nil {
		return nil, normalizeError(err)
	}
	return user, nil
}

func (r *usersRepository) UpdateLastLogin(ctx context.Context, id string) error {
	const query = `UPDATE users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1`
	command, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return normalizeError(err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *usersRepository) UpdateRole(ctx context.Context, id, role string) error {
	command, err := r.db.Exec(ctx, `UPDATE users SET role=$2, updated_at=NOW() WHERE id=$1`, id, role)
	if err != nil {
		return normalizeError(err)
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *usersRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(ctx, `SELECT id, username, email, '', created_at, updated_at, last_login_at, role FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, normalizeError(err)
	}
	defer rows.Close()
	items := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt, &user.Role); err != nil {
			return nil, err
		}
		items = append(items, user)
	}
	return items, rows.Err()
}
