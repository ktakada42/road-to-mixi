package usecase

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/pkg/httputil"
	"problem1/service"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListUseCase interface {
	PostUserLink(ulfr *model.UserLinkForRequest) error
	GetFriendListByUserId(c echo.Context) (*model.FriendList, error)
	GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error)
	GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error)
}

type friendListUseCase struct {
	db  *sql.DB
	fls service.FriendListService
}

func NewFriendListUseCase(db *sql.DB, fls service.FriendListService) FriendListUseCase {
	return &friendListUseCase{
		db:  db,
		fls: fls,
	}
}

func (u *friendListUseCase) checkUserExist(userId int) error {
	exist, err := u.fls.CheckUserExist(userId)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	return httputil.NewHTTPError(err, http.StatusBadRequest, "user not exist")
}

func (u *friendListUseCase) PostUserLink(ulfr *model.UserLinkForRequest) error {
	if err := u.checkUserExist(ulfr.User1Id); err != nil {
		return err
	}
	if err := u.checkUserExist(ulfr.User2Id); err != nil {
		return err
	}

	return u.fls.InsertUserLink(ulfr)
}

func (u *friendListUseCase) GetFriendListByUserId(c echo.Context) (*model.FriendList, error) {
	if err := u.checkUserExist(c.Get("userId").(int)); err != nil {
		return nil, err
	}

	return u.fls.GetFriendListByUserId(c)
}

func (u *friendListUseCase) GetFriendListOfFriendsByUserId(c echo.Context) (*model.FriendList, error) {
	if err := u.checkUserExist(c.Get("userId").(int)); err != nil {
		return nil, err
	}

	return u.fls.GetFriendListOfFriendsByUserId(c)
}

func (u *friendListUseCase) GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) (*model.FriendList, error) {
	if err := u.checkUserExist(c.Get("userId").(int)); err != nil {
		return nil, err
	}

	return u.fls.GetFriendListOfFriendsByUserIdWithPaging(c)
}
