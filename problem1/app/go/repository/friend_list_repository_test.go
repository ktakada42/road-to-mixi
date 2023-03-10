package repository

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"problem1/model"
	"problem1/pkg/testutil"
)

type friendListRepositoryTest struct {
	db  *sql.DB
	flr FriendListRepository
}

func newFriendListRepositoryTest(t *testing.T) *friendListRepositoryTest {
	t.Helper()

	db := testutil.PrepareMySQL(t)
	flr := NewFriendListRepository(db)

	return &friendListRepositoryTest{
		db:  db,
		flr: flr,
	}
}

func newFriendList() *model.FriendList {
	return &model.FriendList{
		Friends: []*model.Friend{
			{
				UserId: 111111,
				Name:   "hoge",
			},
			{
				UserId: 222222,
				Name:   "fuga",
			},
		},
	}
}

type testUser struct {
	userId int
	name   string
}

func newTestUsers() []testUser {
	return []testUser{
		{
			userId: testutil.UserIDForDebug,
			name:   testutil.UserNameForDebug,
		},
		{
			userId: 111111,
			name:   "hoge",
		},
		{
			userId: 222222,
			name:   "fuga",
		},
		{
			userId: 333333,
			name:   "bar",
		},
	}
}

func (r *friendListRepositoryTest) insertTestUserList(t *testing.T, db *sql.DB, tu testUser) {
	t.Helper()

	const q = `
	INSERT INTO users (id, user_id, name)
	VALUES (0, ?, ?)`

	testRecord := []any{
		tu.userId,
		tu.name,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

type userLink struct {
	user1Id int
	user2Id int
}

func newTestUserLink() []userLink {
	return []userLink{
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 111111,
		},
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 222222,
		},
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 333333,
		},
	}

}

func (r *friendListRepositoryTest) insertTestFriendLink(t *testing.T, db *sql.DB, ul userLink) {
	t.Helper()

	const q = `
	INSERT INTO friend_link (id, user1_id, user2_id)
	VALUES (0, ?, ?)`

	testRecord := []any{
		ul.user1Id,
		ul.user2Id,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

func (r *friendListRepositoryTest) insertTestBlockList(t *testing.T, db *sql.DB, ul userLink) {
	t.Helper()

	const q = `
	INSERT INTO block_list (id, user1_id, user2_id)
	VALUES (0, ?, ?)`

	testRecord := []any{
		ul.user1Id,
		ul.user2Id,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

func Test_friendListRepository_InsertUserLink(t *testing.T) {
	tests := []struct {
		name    string
		user1Id int
		user2Id int
		table   string
		want    []int
		wantErr bool
	}{
		{
			name:    "ok: friend_link",
			user1Id: testutil.UserIDForDebug,
			user2Id: 111111,
			table:   "friend_link",
			want:    []int{111111},
			wantErr: false,
		},
		{
			name:    "ok: block_list",
			user1Id: testutil.UserIDForDebug,
			user2Id: 111111,
			table:   "block_list",
			want:    []int{111111},
			wantErr: false,
		},
		{
			name:    "ng: table invalid",
			user1Id: testutil.UserIDForDebug,
			user2Id: 111111,
			table:   "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)

			tx := testutil.BeginTx(t, rt.db)
			err := rt.flr.InsertUserLink(tt.user1Id, tt.user2Id, tt.table)
			if (err != nil) != tt.wantErr {
				testutil.RollBackTx(t, tx)
				t.Fatalf("CheckUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			testutil.CommitTx(t, tx)

			if tt.table == "friend_link" {
				got, err := rt.flr.GetOneHopFriendsUserIdList(tt.user1Id)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.want, got)
			} else {
				got, err := rt.flr.GetBlockUsersIdList(tt.user1Id)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_friendListRepository_CheckUserExist(t *testing.T) {
	userId := testutil.UserIDForDebug
	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				tu := testUser{
					userId: testutil.UserIDForDebug,
					name:   testutil.UserNameForDebug,
				}
				rt.insertTestUserList(t, rt.db, tu)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "ok: user not exist",
			prepare: func(rt *friendListRepositoryTest) {},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.CheckUserExist(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_CheckUserLink(t *testing.T) {
	userLink := newTestUserLink()
	userId := testutil.UserIDForDebug
	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		user2Id int
		table   string
		wantErr bool
	}{
		{
			name: "ok: friend_link",
			prepare: func(rt *friendListRepositoryTest) {
				for _, ul := range userLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			user2Id: 111111,
			table:   "friend_link",
			wantErr: false,
		},
		{
			name: "ok: block_list",
			prepare: func(rt *friendListRepositoryTest) {
				for _, ul := range userLink {
					rt.insertTestBlockList(t, rt.db, ul)
				}
			},
			user2Id: 111111,
			table:   "block_list",
			wantErr: false,
		},
		{
			name:    "ng: friend_link",
			prepare: func(rt *friendListRepositoryTest) {},
			user2Id: 111111,
			table:   "friend_link",
			wantErr: true,
		},
		{
			name:    "ng: block_list",
			prepare: func(rt *friendListRepositoryTest) {},
			user2Id: 111111,
			table:   "block_list",
			wantErr: true,
		},
		{
			name:    "ng: table not exist",
			prepare: func(rt *friendListRepositoryTest) {},
			user2Id: 111111,
			table:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			err := rt.flr.CheckUserLink(userId, tt.user2Id, tt.table)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckUserLink() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func Test_friendListRepository_GetOneHopFriendsUserIdList(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		want    []int
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			want:    []int{111111, 222222, 333333},
			wantErr: false,
		},
		{
			name: "ok: no 1hop friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetOneHopFriendsUserIdList(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetOneHopFrinedsUserIdList() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetBlockUsersIdList(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		want    []int
		wantErr bool
	}{
		{
			name: "ok: 1 user blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
				rt.insertTestBlockList(t, rt.db, userLink{
					user1Id: testutil.UserIDForDebug,
					user2Id: 111111,
				})
			},
			want:    []int{111111},
			wantErr: false,
		},
		{
			name: "ok: all users blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
				for _, ul := range testUserLink {
					rt.insertTestBlockList(t, rt.db, ul)
				}
			},
			want:    []int{111111, 222222, 333333},
			wantErr: false,
		},
		{
			name: "ok: no user blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetBlockUsersIdList(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBlockUsersIdList() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListByUserId(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for i := 0; i < 2; i++ {
					rt.insertTestFriendLink(t, rt.db, testUserLink[i])
				}
			},
			want:    newFriendList(),
			wantErr: false,
		},
		{
			name: "ok: have no friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			want: &model.FriendList{
				Friends: []*model.Friend(nil),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListByUserId(userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListByUserIdExcludingBlockUsers(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name       string
		prepare    func(*friendListRepositoryTest)
		blockUsers []int
		want       *model.FriendList
		wantErr    bool
	}{
		{
			name: "ok: 1 friend blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			blockUsers: []int{333333},
			want:       newFriendList(),
			wantErr:    false,
		},
		{
			name: "ok: all friends blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			blockUsers: []int{111111, 222222, 333333},
			want: &model.FriendList{
				Friends: []*model.Friend(nil),
			},
			wantErr: false,
		},
		{
			name:       "ng: blockUsers nil",
			prepare:    func(rt *friendListRepositoryTest) {},
			blockUsers: nil,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListByUserIdExcludingBlockUsers(userId, tt.blockUsers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserIdExcludingBlockUsers() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListOfFriendsByUserId(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := []userLink{
		{
			user1Id: 123456789,
			user2Id: 444444,
		},
		{
			user1Id: 444444,
			user2Id: 111111,
		},
		{
			user1Id: 444444,
			user2Id: 222222,
		},
		{
			user1Id: 444444,
			user2Id: 333333,
		},
	}
	testUserLink2 := newTestUserLink()

	tests := []struct {
		name         string
		prepare      func(*friendListRepositoryTest)
		excludeUsers []int
		want         *model.FriendList
		wantErr      bool
	}{
		{
			name: "ok: 1 friend excluded",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			excludeUsers: []int{333333},
			want:         newFriendList(),
			wantErr:      false,
		},
		{
			name: "ok: all friends excluded",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			excludeUsers: []int{111111, 222222, 333333},
			want: &model.FriendList{
				Friends: []*model.Friend(nil),
			},
			wantErr: false,
		},
		{
			name: "ok: have no 2hop friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink2 {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			excludeUsers: []int{444444},
			want: &model.FriendList{
				Friends: []*model.Friend(nil),
			},
			wantErr: false,
		},
		{
			name: "ok: have no friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			excludeUsers: []int{111111},
			want: &model.FriendList{
				Friends: []*model.Friend(nil),
			},
			wantErr: false,
		},
		{
			name:         "ng: excludeUsers nil",
			prepare:      func(rt *friendListRepositoryTest) {},
			excludeUsers: nil,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListOfFriendsByUserId(userId, tt.excludeUsers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListOfFriendsByUserIdWithPaging(t *testing.T) {
	userId := testutil.UserIDForDebug
	testUsers := newTestUsers()
	testUserLink := []userLink{
		{
			user1Id: 123456789,
			user2Id: 444444,
		},
		{
			user1Id: 444444,
			user2Id: 111111,
		},
		{
			user1Id: 444444,
			user2Id: 222222,
		},
		{
			user1Id: 444444,
			user2Id: 333333,
		},
	}

	tests := []struct {
		name         string
		prepare      func(*friendListRepositoryTest)
		excludeUsers []int
		limit        int
		offset       int
		want         *model.FriendList
		wantErr      bool
	}{
		{
			name: "ok: limit",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			excludeUsers: []int{444444},
			limit:        2,
			offset:       0,
			want:         newFriendList(),
			wantErr:      false,
		},
		{
			name: "ok: offset",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			excludeUsers: []int{444444},
			limit:        3,
			offset:       2,
			want: &model.FriendList{
				Friends: []*model.Friend{
					{
						UserId: 333333,
						Name:   "bar",
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "ng: excludeUsers nil",
			prepare:      func(rt *friendListRepositoryTest) {},
			excludeUsers: nil,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListOfFriendsByUserIdWithPaging(userId, tt.excludeUsers, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserIdWithPaging() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
