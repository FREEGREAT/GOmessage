package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/postgresql"
)

type userRepository struct {
	client postgresql.Client
}

// Create implements user.UserRepository.
func (r *userRepository) Create(ctx context.Context, user *models.UserModel) error {
	q := `INSERT INTO users
			(nickname, password_hash, email, age, image_url) 
		VALUES 
			($1, $2, $3, $4,$5)
		RETURNING user_id, nickname
	`
	row := r.client.QueryRow(ctx, q, user.Nickname, user.PasswordHash, user.Email, user.Age, user.ImageUrl)

	err := row.Scan(&user.ID, &user.Nickname)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Sprintf("Sql Error: %s, Detail %s, Where %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			return err
		}
		return err
	}
	return nil
}

// Delete implements user.UserRepository.
func (r *userRepository) Delete(ctx context.Context, id string) (string, error) {
	var nickname string
	query := `DELETE FROM users WHERE user_id = $1 RETURNING nickname`
	row := r.client.QueryRow(ctx, query, id)
	err := row.Scan(&nickname)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Sprintf("Sql Error: %s, Detail %s, Where %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			return "", err
		}
		return "", err
	}
	return nickname, nil
}

// FindAll implements user.UserRepository.
func (r *userRepository) FindAll(ctx context.Context) (u []models.UserModel, err error) {
	querry := `SELECT user_id, nickname, password_hash, email, age, image_url FROM users`
	qrow, err := r.client.Query(ctx, querry)
	if err != nil {
		return nil, err
	}

	users := make([]models.UserModel, 0)
	for qrow.Next() {
		var usr models.UserModel
		err := qrow.Scan(&usr.ID, &usr.Nickname, &usr.PasswordHash, &usr.Email, &usr.Age, &usr.ImageUrl)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	if err = qrow.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// FindOne implements user.UserRepository.
func (r *userRepository) FindOne(ctx context.Context, id string) (models.UserModel, error) {
	querry := `SELECT user_id, nickname, password_hash, email, age, image_url FROM users WHERE user_id = $1`
	qrow := r.client.QueryRow(ctx, querry, id)
	var usr models.UserModel
	if err := qrow.Scan(&usr.ID, &usr.Nickname, &usr.PasswordHash, &usr.Email, &usr.Age, &usr.ImageUrl); err != nil {
		return models.UserModel{}, err
	}
	return usr, nil
}

// Update implements user.UserRepository.
func (r *userRepository) Update(ctx context.Context, users *models.UserModel) error {
	query := "UPDATE users SET"
	var args []interface{}
	var setClauses []string
	argIdx := 1

	if users.Nickname != "" {
		setClauses = append(setClauses, "nickname = $"+fmt.Sprint(argIdx))
		args = append(args, users.Nickname)
		argIdx++
	}
	if users.PasswordHash != "" {
		setClauses = append(setClauses, "password_hash = $"+fmt.Sprint(argIdx))
		args = append(args, users.PasswordHash)
		argIdx++
	}
	if users.Email != "" {
		setClauses = append(setClauses, "email = $"+fmt.Sprint(argIdx))
		args = append(args, users.Email)
		argIdx++
	}
	if users.Age != nil {
		setClauses = append(setClauses, "age = $"+fmt.Sprint(argIdx))
		args = append(args, users.Age)
		argIdx++
	}
	if users.ImageUrl != nil {
		setClauses = append(setClauses, "image_url = $"+fmt.Sprint(argIdx))
		args = append(args, users.ImageUrl)
		argIdx++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += " " + strings.Join(setClauses, ", ") + " WHERE user_id = $" + fmt.Sprint(argIdx)
	args = append(args, users.ID)

	cmdTag, err := r.client.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows were updated for user_id: %s", users.ID)
	}

	return nil
}

func NewUserRepository(client postgresql.Client) storage.UserRepository {
	return &userRepository{
		client: client,
	}
}
