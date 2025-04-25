package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testLicenseModels models.LicenseModels

func mockLicenseRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "key", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockLicense())...)

}

func mockLicense() models.License {
	return models.License{
		ID:        "1234567890",
		Key:       "test_key",
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func TestGetLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"licenses\" WHERE id = \\$1").
		WithArgs("1234567890").
		WillReturnRows(mockLicenseRows())

	models.MockSetup(db)
	license, err := testLicenseModels.Get("1234567890")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := mockLicense()
	assert.Equal(t, expected, license, "they should be equal")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateLicense(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"licenses\"").
		WithArgs(StructToSlice(mockLicense())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	lincese := mockLicense()
	err = testLicenseModels.Create(&lincese)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
