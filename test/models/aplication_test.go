package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testApplicationModels models.ApplicationModels

func mockApplicationRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "description", "secret_key", "salt", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockApplication0())...)
}

func mockApplicationsRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "description", "secret_key", "salt", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockApplication0())...).
		AddRow(StructToSlice(mockApplication1())...)
}

func mockApplication0() models.Application {
	return models.Application{
		ID:          "001",
		Name:        "test_app1",
		Description: "test_app1 description",
		SecretKey:   "secret_key1",
		Salt:        "salt1",
		CreatedAt:   time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:      true,
	}
}

func mockApplication1() models.Application {
	return models.Application{
		ID:          "002",
		Name:        "test_app2",
		Description: "test_app2 description",
		SecretKey:   "secret_key2",
		Salt:        "salt2",
		CreatedAt:   time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:      true,
	}
}

func TestGetApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"applications\" WHERE id = \\$1").
		WithArgs("001").
		WillReturnRows(mockApplicationRows())

	models.MockSetup(db)
	app, err := testApplicationModels.Get("001")
	if err != nil {
		t.Errorf("error was not expected while getting application: %s", err)
	}

	expect := mockApplication0()
	assert.Equal(t, expect, app, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApplicationCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"applications\" WHERE status = true").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	models.MockSetup(db)
	count, err := testApplicationModels.Count()
	if err != nil {
		t.Errorf("error was not expected while counting applications: %s", err)
	}

	expect := int64(2)
	assert.Equal(t, expect, count, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListApplications(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"applications\" WHERE status = true ORDER BY created_at desc LIMIT \\$1 OFFSET \\$2").
		WithArgs(10, 1).
		WillReturnRows(mockApplicationsRows())

	models.MockSetup(db)
	apps, err := testApplicationModels.List(10, 1, "created_at desc")
	if err != nil {
		t.Errorf("error was not expected while listing applications: %s", err)
	}

	expect := []models.Application{mockApplication0(), mockApplication1()}
	assert.Equal(t, expect, apps, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	application := mockApplication0()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"applications\"").
		WithArgs(StructToSlice(application)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testApplicationModels.Create(&application)
	if err != nil {
		t.Errorf("error was not expected while creating application: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateApplication(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	application := mockApplication0()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"applications\"").
		WithArgs(StructToUpdateSlice(application)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testApplicationModels.Update(&application)
	if err != nil {
		t.Errorf("error was not expected while updating application: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
