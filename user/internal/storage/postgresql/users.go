package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/postgresql"
)

type repository struct {
	client postgresql.Client
}

// Create implements user.UserRepository.
func (r *repository) Create(ctx context.Context, user *models.UserModel) error {
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
func (r *repository) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements user.UserRepository.
func (r *repository) FindAll(ctx context.Context) (u []models.UserModel, err error) {
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
func (r *repository) FindOne(ctx context.Context, id string) (models.UserModel, error) {
	querry := `SELECT user_id, nickname, password_hash, email, age, image_url * FROM users WHERE user_id = $1`
	qrow := r.client.QueryRow(ctx, querry, id)
	var usr models.UserModel
	if err := qrow.Scan(&usr.ID, &usr.Nickname, &usr.PasswordHash, &usr.Email, &usr.Age, &usr.ImageUrl); err != nil {
		return models.UserModel{}, err
	}
	return usr, nil
}

// Update implements user.UserRepository.
func (r *repository) Update(ctx context.Context, users models.UserModel) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client) storage.UserRepository {
	return &repository{
		client: client,
	}
}
