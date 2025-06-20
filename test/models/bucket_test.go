package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testBucketModels models.BucketModels

func mockBucketRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "region", "provider", "endpoint", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockBucket0())...)
}

func mockBucketsRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "region", "provider", "endpoint", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockBucket0())...).
		AddRow(StructToSlice(mockBucket1())...)
}

func mockBucket0() models.Bucket {
	return models.Bucket{
		ID:        "test_id",
		Name:      "test_name",
		Region:    "us-west-2",
		Provider:  "aws",
		Endpoint:  "s3.us-west-2.amazonaws.com",
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func mockBucket1() models.Bucket {
	return models.Bucket{
		ID:        "test_id_1",
		Name:      "test_name_1",
		Region:    "us-west-2",
		Provider:  "aws",
		Endpoint:  "s3.us-west-2.amazonaws.com",
		CreatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:    true,
	}
}

func TestGetBucket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"buckets\" WHERE id = \\$1").
		WithArgs("test_id").
		WillReturnRows(mockBucketRows())

	models.MockSetup(db)
	bucket, err := testBucketModels.Get("test_id")
	assert.Equal(t, err, nil)
	executed := mockBucket0()
	assert.Equal(t, executed, bucket)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBucketByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"buckets\" WHERE name = \\$1 and status = true").
		WithArgs("test_name").
		WillReturnRows(mockBucketRows())

	models.MockSetup(db)
	bucket, err := testBucketModels.GetByName("test_name")
	assert.Equal(t, err, nil)

	executed := mockBucket0()
	assert.Equal(t, executed, bucket)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountBuckets(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"buckets\" WHERE status = true").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	models.MockSetup(db)
	count, err := testBucketModels.Count()
	assert.Equal(t, err, nil)
	executed := int64(2)
	assert.Equal(t, executed, count)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListBuckets(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"buckets\" WHERE status = true ORDER BY created_at desc LIMIT \\$1 OFFSET \\$2").
		WithArgs(10, 1).
		WillReturnRows(mockBucketsRows())

	models.MockSetup(db)
	buckets, err := testBucketModels.List(10, 1, "created_at desc")
	assert.Equal(t, err, nil)

	executed := []models.Bucket{mockBucket0(), mockBucket1()}
	assert.Equal(t, executed, buckets)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateBucket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	bucket := mockBucket0()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"buckets\"").
		WithArgs(StructToSlice(mockBucket0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testBucketModels.Create(&bucket)
	assert.Equal(t, err, nil)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateBucket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	bucket := mockBucket0()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"buckets\" SET").
		WithArgs(StructToUpdateSlice(mockBucket0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testBucketModels.Update(&bucket)
	assert.Equal(t, err, nil)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteBucket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"buckets\" WHERE id = \\$1").
		WithArgs("test_id").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testBucketModels.Delete("test_id")
	assert.Equal(t, err, nil)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
