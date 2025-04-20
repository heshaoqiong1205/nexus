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
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery("SELECT \\* FROM \"licenses\" WHERE id = \\$1").
		WithArgs("1234567890").
		WillReturnRows(mockLicenseRows())

	models.MockSetup(mockDB)
	license, err := testLicenseModels.Get("1234567890")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := mockLicense()
	assert.Equal(t, expected, license, "they should be equal")
}

func TestCreateLicense(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"licenses\"").
		WithArgs(StructToSlice(mockLicense())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(mockDB)
	lincese := mockLicense()
	err = testLicenseModels.Create(&lincese)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
