package postgresql

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/postgresql"
)

type friendsRepository struct {
	client postgresql.Client
}

const nil_string = " "

// Create implements storage.FriendsRepository.
func (f *friendsRepository) Create(ctx context.Context, friend *models.FriendListModel) (string, error) {
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
			var info_err string
			if pgErr.Code == pgerrcode.UniqueViolation {
				info_err = "You already friends"
			} else if pgErr.Code == pgerrcode.InvalidTextRepresentation {
				info_err = "Invalid friend id"
			}
			return info_err, err
		}
		return nil_string, err
	}
	return nil_string, nil
}

func (f *friendsRepository) FindAll(ctx context.Context, user_id string) ([]models.FriendListModel, error) {
	q := `SELECT friend_id FROM friend_list WHERE user_id=$1`
	qrow, err := f.client.Query(ctx, q, user_id)
	if err != nil {
		return nil, err
	}

	friends := make([]models.FriendListModel, 0)
	for qrow.Next() {
		var frn models.FriendListModel
		err := qrow.Scan(&frn.FriendID)
		if err != nil {
			return nil, err
		}
		friends = append(friends, frn)
	}
	if err = qrow.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

// Delete implements storage.FriendsRepository.
func (f *friendsRepository) Delete(ctx context.Context, friends *models.FriendListModel) error {
	query := `DELETE FROM friend_list WHERE user_id = $1 AND friend_id = $2 OR friend_id=$1 AND user_id =$2`
	_, err := f.client.Exec(ctx, query, friends.UserID, friends.FriendID)
	if err != nil {
		return err
	}
	return nil
}

func NewFriendsRepository(client postgresql.Client) storage.FriendsRepository {
	return &friendsRepository{
		client: client,
	}
}
