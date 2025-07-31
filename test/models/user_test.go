package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testUserModels models.UserModels

func mockUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "account", "username", "password", "region",
		"location", "icon", "role", "last_login_time", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockUser0())...)
}

func mockUsersRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "account", "username", "password", "region",
		"location", "icon", "role", "last_login_time", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockUser0())...).
		AddRow(StructToSlice(mockUser1())...)
}

func mockUser0() models.User {
	location := models.NewPoint(121.123456, 31.123456)

	return models.User{
		ID:            "001",
		Account:       "account1",
		Username:      "test_user1",
		Password:      "password1",
		Region:        "region1",
		Location:       location,
		Icon:          "icon1",
		Role:          "admin",
		LastLoginTime: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		CreatedAt:     time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:        true,
	}
}

func mockUser1() models.User {
	location := models.NewPoint(34.0522, -118.2437)

	return models.User{
		ID:            "002",
		Account:       "account2",
		Username:      "test_user2",
		Password:      "password2",
		Region:        "region2",
		Location:      location,
		Icon:          "icon2",
		Role:          "user",
		LastLoginTime: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		CreatedAt:     time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:        true,
	}
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs("001").
		WillReturnRows(mockUserRows())

	models.MockSetup(db)
	user, err := testUserModels.Get("001")
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	expect := mockUser0()
	assert.Equal(t, expect, user, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE account = \\$1").
		WithArgs("account1").
		WillReturnRows(mockUserRows())

	models.MockSetup(db)
	user, err := testUserModels.GetByAccount("account1")
	if err != nil {
		t.Errorf("error was not expected while getting user by account: %s", err)
	}

	expect := mockUser0()
	assert.Equal(t, expect, user, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	likeAccount := "account%"
	mock.ExpectQuery("SELECT \\* FROM \"users\" WHERE status = true and account like \\$1 ORDER BY created_at desc LIMIT \\$2 OFFSET \\$3").
		WithArgs(likeAccount, 10, 1).
		WillReturnRows(mockUsersRows())

	models.MockSetup(db)
	users, err := testUserModels.List(likeAccount, 10, 1, "created_at desc")
	if err != nil {
		t.Errorf("error was not expected while listing users: %s", err)
	}

	expect := []models.User{mockUser0(), mockUser1()}
	assert.Equal(t, expect, users, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	user := mockUser0()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"users\"").
		WithArgs(StructToSlice(mockUser0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testUserModels.Create(&user)
	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	user := mockUser0()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"users\" SET").
		WithArgs(StructToUpdateSlice(mockUser0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testUserModels.Update(&user)
	if err != nil {
		t.Errorf("error was not expected while updating user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
