package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/postgresql"
)

type friendsRepository struct {
	client postgresql.Client
}

// Create implements storage.FriendsRepository.
func (f *friendsRepository) Create(ctx context.Context, friend *models.FriendListModel) error {
	q := `INSERT INTO friend_list
			(user_id, friend_id) 
		VALUES 
			($1, $2)
		RETURNING id
`
	row := f.client.QueryRow(ctx, q, friend.UserID, friend.FriendID)

	err := row.Scan(&friend.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Sprintf("Sql Error: %s, Detail %s, Where %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			return err
		}
		return err
	}
	return nil
}

// Delete implements storage.FriendsRepository.
func (f *friendsRepository) Delete(ctx context.Context, friends *models.FriendListModel) error {
	query := `DELETE FROM friend_list WHERE user_id = $1 AND friend_id = $2`
	f.client.QueryRow(ctx, query, friends.UserID, friends.FriendID)

	return nil
}

func NewFriendsRepository(client postgresql.Client) storage.FriendsRepository {
	return &friendsRepository{
		client: client,
	}
}
