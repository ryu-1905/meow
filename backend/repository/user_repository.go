package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository định nghĩa các thao tác với bảng users trong cơ sở dữ liệu.
type UserRepository interface {
	GetHashPasswordFromEmail(ctx context.Context, email string) (int64, string, error)
	CreateUser(ctx context.Context, name, email, hashPassword string) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository khởi tạo một UserRepository mới.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetHashPasswordFromEmail lấy id và hash_password của người dùng dựa trên email.
func (r *userRepository) GetHashPasswordFromEmail(ctx context.Context, email string) (int64, string, error) {
	const query = "SELECT id, hash_password FROM users WHERE email = $1"

	slog.Info("Executing SQL query", slog.String("query", query))
	slog.Debug("Query parameters", slog.String("email", email))

	var id int64
	var hashPassword string

	err := r.db.QueryRow(ctx, query, email).Scan(&id, &hashPassword)
	if err != nil {
		slog.Error("SQL execution failed", slog.String("query", query), slog.Any("error", err))
		slog.Debug("SQL execution failed details", slog.String("email", email))
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", fmt.Errorf("user not found with email: %s", email)
		}
		return 0, "", fmt.Errorf("failed to get hash password from email: %w", err)
	}

	slog.Info("SQL execution succeeded")
	slog.Debug("SQL execution succeeded details", slog.Int64("user_id", id))
	return id, hashPassword, nil
}

// CreateUser tạo một người dùng mới trong cơ sở dữ liệu và trả về ID của bản ghi vừa tạo.
func (r *userRepository) CreateUser(ctx context.Context, name, email, hashPassword string) (int64, error) {
	const query = "INSERT INTO users (name, email, hash_password) VALUES ($1, $2, $3) RETURNING id"

	slog.Info("Executing SQL query", slog.String("query", query))
	slog.Debug("Query parameters", slog.String("name", name), slog.String("email", email))

	var id int64
	err := r.db.QueryRow(ctx, query, name, email, hashPassword).Scan(&id)
	if err != nil {
		slog.Error("SQL execution failed", slog.String("query", query), slog.Any("error", err))
		slog.Debug("SQL execution failed details", slog.String("name", name), slog.String("email", email))
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("SQL execution succeeded")
	slog.Debug("SQL execution succeeded details", slog.Int64("created_user_id", id))
	return id, nil
}

// DeleteUser xóa một người dùng dựa trên ID.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	const query = "DELETE FROM users WHERE id = $1"

	slog.Info("Executing SQL query", slog.String("query", query))
	slog.Debug("Query parameters", slog.Int64("id", id))

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("SQL execution failed", slog.String("query", query), slog.Any("error", err))
		slog.Debug("SQL execution failed details", slog.Int64("id", id))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	slog.Info("SQL execution succeeded")
	return nil
}
