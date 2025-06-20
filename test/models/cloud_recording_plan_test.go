package models_testing

import (
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testCloudRecordingPlanModels models.CloudRecordingPlanModels

func mockCloudRecordingPlanRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "channel", "enabled", "mode", "mfd", "interval", "storage_id", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudRecordingPlan0())...)
}

func mockCloudRecordingPlansRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "device_id", "channel", "enabled", "mode", "mfd", "interval", "storage_id", "created_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockCloudRecordingPlan0())...).
		AddRow(StructToSlice(mockCloudRecordingPlan1())...)
}

func mockCloudRecordingPlan0() models.CloudRecordingPlan {
	return models.CloudRecordingPlan{
		ID:        "plan_id_1",
		DeviceID:  "device_id_1",
		Channel:   0,
		Enabled:   true,
		Mode:      "event",
		MFD:       60,
		Interval:  30,
		StorageID: "storage_id_1",
		CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		Status:    true,
	}
}

func mockCloudRecordingPlan1() models.CloudRecordingPlan {
	return models.CloudRecordingPlan{
		ID:        "plan_id_2",
		DeviceID:  "device_id_1",
		Channel:   1,
		Enabled:   false,
		Mode:      "continuous",
		MFD:       120,
		Interval:  60,
		StorageID: "storage_id_2",
		CreatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
		Status:    true,
	}
}

func TestGetCloudRecordingPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_recording_plans\" WHERE id = \\$1").
		WithArgs("plan_id_1").
		WillReturnRows(mockCloudRecordingPlanRows())
	models.MockSetup(db)

	plan, err := testCloudRecordingPlanModels.Get("plan_id_1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := mockCloudRecordingPlan0()
	assert.Equal(t, expected, plan)
}

func TestGetCloudRecordingPlansByDevice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_recording_plans\" WHERE device_id = \\$1 and status = true").
		WithArgs("device_id_1").
		WillReturnRows(mockCloudRecordingPlansRows())
	models.MockSetup(db)

	plans, err := testCloudRecordingPlanModels.GetByDevice("device_id_1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := []models.CloudRecordingPlan{mockCloudRecordingPlan0(), mockCloudRecordingPlan1()}
	assert.Equal(t, expected, plans)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCountCloudRecordingPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"cloud_recording_plans\" WHERE status = true").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	models.MockSetup(db)

	count, err := testCloudRecordingPlanModels.Count()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	assert.Equal(t, int64(2), count)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListCloudRecordingPlans(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"cloud_recording_plans\" WHERE status = true ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
		WithArgs(10, 1).
		WillReturnRows(mockCloudRecordingPlansRows())
	models.MockSetup(db)

	plans, err := testCloudRecordingPlanModels.List(10, 1, "created_at DESC")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := []models.CloudRecordingPlan{mockCloudRecordingPlan0(), mockCloudRecordingPlan1()}
	assert.Equal(t, expected, plans)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateCloudRecordingPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	plan := mockCloudRecordingPlan0()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"cloud_recording_plans\"").
		WithArgs(StructToSlice(plan)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudRecordingPlanModels.Create(&plan)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCloudRecordingPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	plan := mockCloudRecordingPlan0()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"cloud_recording_plans\" SET").
		WithArgs(StructToUpdateSlice(plan)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudRecordingPlanModels.Update(&plan)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCloudRecordingPlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"cloud_recording_plans\" WHERE id = \\$1").
		WithArgs("plan_id_1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testCloudRecordingPlanModels.Delete("plan_id_1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
