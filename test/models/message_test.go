package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testMessageModels models.MessageModels

func mockMessagesRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "user_id", "type", "content", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockMessage0())...).
		AddRow(StructToSlice(mockMessage1())...)
}

func mockMessage0() models.Message {
	return models.Message{
		ID:        "test_id",
		UserID:    "test_user_id",
		Type:      "test_type",
		Content:   []byte("test_content"),
		CreatedAt: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		Status:    true,
	}
}

func mockMessage1() models.Message {
	return models.Message{
		ID:        "test_id_1",
		UserID:    "test_user_id_1",
		Type:      "test_type_1",
		Content:   []byte("test_content_1"),
		CreatedAt: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		Status:    true,
	}
}

func TestCreateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	message := mockMessage0()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"messages\"").
		WithArgs(StructToSlice(message)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testMessageModels.Create(&message)
	if err != nil {
		t.Errorf("error was not expected while creating message: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	message := mockMessage0()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"messages\" SET").
		WithArgs(StructToUpdateSlice(mockMessage0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testMessageModels.Update(&message)
	if err != nil {
		t.Errorf("error was not expected while updating message: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	userID := "test_user_id"

	mock.ExpectQuery("SELECT \\* FROM \"messages\"").
		WillReturnRows(mockMessagesRows())

	models.MockSetup(db)
	messages, err := testMessageModels.List(&userID, nil, 10, 0)
	if err != nil {
		t.Errorf("error was not expected while listing messages: %s", err)
	}

	expect := []models.Message{mockMessage0(), mockMessage1()}
	assert.Equal(t, expect, messages, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
