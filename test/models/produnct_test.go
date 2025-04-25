package models_testing

import (
	"encoding/json"
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testProductModels models.ProductModels

func mockProductRows(product models.Product) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "description", "required_features", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(product)...)

}

func mockProduct() models.Product {
	required_features, _ := json.Marshal([]string{"p2p", "cloud_storage"})
	// time_now := time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC)
	return models.Product{
		ID:               "s1",
		Name:             "Test Product",
		Description:      "This is a test product",
		RequiredFeatures: required_features,
		CreatedAt:        time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2024, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:           true,
	}
}

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := mockProduct()
	rows := mockProductRows(product)
	mock.ExpectQuery("SELECT \\* FROM \"products\" WHERE id = \\$1").
		WithArgs(product.ID).
		WillReturnRows(rows)

	models.MockSetup(db)
	result, err := testProductModels.Get(product.ID)
	if err != nil {
		t.Errorf("Error getting product: %v", err)
	}
	assert.Equal(t, product, result, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := mockProduct()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"products\"").
		WithArgs(StructToSlice(product)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	var featrues []string
	if json.Unmarshal(product.RequiredFeatures, &featrues) != nil {
		t.Errorf("Error unmarshalling required features: %v", err)
	}
	err = testProductModels.Create(&product)
	if err != nil {
		t.Errorf("Expected a non-nil product: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	product := mockProduct()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"products\" SET").
		WithArgs(StructToUpdateSlice(product)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testProductModels.Update(&product)
	if err != nil {
		t.Errorf("Error updating product: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
