package repository

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v4"

	"problem1/model"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListRepository interface {
	CheckUserExist(c echo.Context) (bool, error)
	GetFriendListByUserId(c echo.Context) ([]*model.User, error)
}

type friendListRepository struct {
	db *sql.DB
}

func NewFriendListRepository(db *sql.DB) FriendListRepository {
	return &friendListRepository{
		db: db,
	}
}

func (r *friendListRepository) CheckUserExist(c echo.Context) (bool, error) {
	userId := c.QueryParam("userId")

	const q = `
SELECT user_id, name
FROM users
WHERE user_id = ?`

	row := r.db.QueryRow(q, userId)

	user := &model.User{}
	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *friendListRepository) GetFriendListByUserId(c echo.Context) ([]*model.User, error) {
	userId := c.QueryParam("userId")

	const q = `
SELECT U.user_id, U.name
FROM users AS U INNER JOIN friend_link AS FL
ON U.user_id = FL.user2_id
WHERE FL.user1_id = ?`

	rows, err := r.db.Query(q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendList []*model.User
	for rows.Next() {
		friend := &model.User{}
		if err := rows.Scan(&friend.Id, &friend.Name); err != nil {
			return nil, err
		}

		friendList = append(friendList, friend)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friendList, nil
}