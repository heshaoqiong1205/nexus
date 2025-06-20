package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testCloudStorage models.CloudStorageModels

func mockCloudStorageRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "tos", "bucket", "mode", "path", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudStorage0())...)
}

func mockCloudStoragesRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "tos", "bucket", "mode", "path", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudStorage0())...).
		AddRow(StructToSlice(mockCloudStorage1())...)
}

func mockCloudStorage0() models.CloudStorage {
	return models.CloudStorage{
		ID:        "1",
		DeviceID:  "device1",
		Tos:       "log",
		Bucket:    "bucket_name1",
		Mode:      "encrypted",
		Path:      "/88012e-34125598-194c9d2a7048680e/media",
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func mockCloudStorage1() models.CloudStorage {
	return models.CloudStorage{
		ID:        "2",
		DeviceID:  "device1",
		Tos:       "media30",
		Bucket:    "bucket_name2",
		Mode:      "plaintext",
		Path:      "/88012e-34125598-194c9d2a7048680e/media",
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func TestGetCloudStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_storages\" WHERE id = \\$1").
		WithArgs("1").
		WillReturnRows(mockCloudStorageRows())

	models.MockSetup(db)
	storage, err := testCloudStorage.Get("1")
	if err != nil {
		t.Errorf("error was not expected while getting cloud storage: %s", err)
	}

	executed := mockCloudStorage0()
	assert.Equal(t, executed, storage)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCloudStorageByDeviceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_storages\" WHERE device_id = \\$1 and status = true").
		WithArgs("device1").
		WillReturnRows(mockCloudStoragesRows())

	models.MockSetup(db)
	storages, err := testCloudStorage.GetByDeviceID("device1")
	if err != nil {
		t.Errorf("error was not expected while getting cloud storage by device id: %s", err)
	}

	executed := []models.CloudStorage{mockCloudStorage0(), mockCloudStorage1()}
	assert.Equal(t, executed, storages)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCloudStorageList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_storages\" WHERE status = true ORDER BY created_at desc LIMIT \\$1 OFFSET \\$2").
		WithArgs(10, 1).
		WillReturnRows(mockCloudStoragesRows())

	models.MockSetup(db)
	storages, err := testCloudStorage.List(10, 1, "created_at desc")
	if err != nil {
		t.Errorf("error was not expected while getting cloud storage count: %s", err)
	}

	executed := []models.CloudStorage{mockCloudStorage0(), mockCloudStorage1()}
	assert.Equal(t, executed, storages)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateCloudStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"cloud_storages\"").
		WithArgs(StructToSlice(mockCloudStorage0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	cloudStorage := mockCloudStorage0()
	err = testCloudStorage.Create(&cloudStorage)
	if err != nil {
		t.Errorf("error was not expected while creating cloud storage: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCloudStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"cloud_storages\" SET").
		WithArgs(StructToUpdateSlice(mockCloudStorage0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	cloudStorage := mockCloudStorage0()
	err = testCloudStorage.Update(&cloudStorage)
	if err != nil {
		t.Errorf("error was not expected while updating cloud storage: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCloudStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"cloud_storages\" WHERE id = \\$1").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudStorage.Delete("1")
	if err != nil {
		t.Errorf("error was not expected while deleting cloud storage: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
