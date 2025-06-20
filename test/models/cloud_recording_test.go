package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testCloudRecordingModels models.CloudRecordingModels

func mockCloudRecordingRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "channel", "fragment_duration", "bucket_id", "prefix", "fragments", "state", "begin_time", "end_time", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudRecording0())...)
}

func mockCloudRecordingsRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "channel", "fragment_duration", "bucket_id", "prefix", "fragments", "state", "begin_time", "end_time", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudRecording0())...).
		AddRow(StructToSlice(mockCloudRecording1())...)
}

func mockCloudRecording0() models.CloudRecording {
	return models.CloudRecording{
		ID:               "test_id",
		DeviceID:         "test_device_id",
		Channel:          0,
		FragmentDuration: 10,
		BucketID:         "test_bucket_id",
		Prefix:           "test_prefix",
		Fragments:        []byte(`[0, 10, 20]`),
		State:            "recording",
		BeginTime:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		EndTime:          time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC),
		CreatedAt:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		Status:           true,
	}
}

func mockCloudRecording1() models.CloudRecording {
	return models.CloudRecording{
		ID:               "test_id_1",
		DeviceID:         "test_device_id",
		Channel:          0,
		FragmentDuration: 10,
		BucketID:         "test_bucket_id_1",
		Prefix:           "test_prefix_1",
		Fragments:        []byte(`[0, 10, 20, 30]`),
		State:            "completed",
		BeginTime:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		EndTime:          time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC),
		CreatedAt:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC),
		Status:           true,
	}
}

func TestGetCloudRecording(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_recordings\" WHERE id = \\$1").
		WithArgs("test_id").
		WillReturnRows(mockCloudRecordingRows())

	models.MockSetup(db)
	recording, err := testCloudRecordingModels.Get("test_id")
	if err != nil {
		t.Errorf("error was not expected while getting cloud recording: %s", err)
	}
	assert.Equal(t, mockCloudRecording0(), recording)
}

func TestCountCloudRecordings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"cloud_recordings\" WHERE device_id = \\$1 and channel = \\$2 and begin_time >= \\$3 and begin_time =< \\$4 and status = true").
		WithArgs("test_device_id", 0, time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	models.MockSetup(db)
	count, err := testCloudRecordingModels.Count("test_device_id", 0, time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC))
	if err != nil {
		t.Errorf("error was not expected while getting cloud recording count: %s", err)
	}
	assert.Equal(t, int64(2), count)
}

func TestListCloudRecordings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_recordings\" WHERE device_id = \\$1 and channel = \\$2 and begin_time >= \\$3 and begin_time =< \\$4 and status = true ORDER BY begin_time desc LIMIT \\$5 OFFSET \\$6").
		WithArgs("test_device_id", 0, time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC), 10, 1).
		WillReturnRows(mockCloudRecordingsRows())

	models.MockSetup(db)
	recordings, err := testCloudRecordingModels.List("test_device_id", 0, time.Date(2023, time.October, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, time.October, 1, 1, 0, 0, 0, time.UTC), 10, 1)
	if err != nil {
		t.Errorf("error was not expected while getting cloud recording list: %s", err)
	}
	executed := []models.CloudRecording{mockCloudRecording0(), mockCloudRecording1()}
	assert.Equal(t, executed, recordings)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestCreateCloudRecording(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	recording := mockCloudRecording0()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"cloud_recordings\"").
		WithArgs(StructToSlice(mockCloudRecording0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudRecordingModels.Create(&recording)
	if err != nil {
		t.Errorf("error was not expected while creating cloud recording: %s", err)
	}
}

func TestUpdateCloudRecording(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	recording := mockCloudRecording0()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"cloud_recordings\" SET").
		WithArgs(StructToUpdateSlice(mockCloudRecording0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudRecordingModels.Update(&recording)
	if err != nil {
		t.Errorf("error was not expected while updating cloud recording: %s", err)
	}
}
