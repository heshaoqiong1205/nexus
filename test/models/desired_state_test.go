package models_testing

import (
	"encoding/json"
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testDesiredStateModels models.DesiredStateModels

func mockDesiredStateRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "state", "version", "status", "last_desired_id", "confirmed_at", "created_at", "updated_at"}).
		AddRow(StructToSlice(mockDesiredState0())...)
}

func mockDesiredStatesRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "state", "version", "status", "last_desired_id", "confirmed_at", "created_at", "updated_at"}).
		AddRow(StructToSlice(mockDesiredState0())...).
		AddRow(StructToSlice(mockDesiredState1())...)
}

func mockDesiredState0() models.DesiredState {
	desiredStateJSON, _ := json.Marshal(map[string]interface{}{
		"volume": map[string]interface{}{
			"value":     85,
			"timestamp": time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":        1234,
		},
		"record": map[string]interface{}{
			"status":    true,
			"mode":      2,
			"timestamp": time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":        1235,
		},
		"video": map[string]interface{}{
			"brightness": 70,
			"contrast":   60,
			"timestamp":  time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":         1236,
		},
	})

	return models.DesiredState{
		ID:            "dev-001", // DeviceID as primary key
		State:         desiredStateJSON,
		Version:       1,
		Status:        true,
		LastDesiredID: 1001,
		ConfirmedAt:   nil,
		CreatedAt:     time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC),
	}
}

func mockDesiredState1() models.DesiredState {
	desiredStateJSON, _ := json.Marshal(map[string]interface{}{
		"volume": map[string]interface{}{
			"value":     90,
			"mute":      false,
			"timestamp": time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":        1234,
		},
		"record": map[string]interface{}{
			"status":    false,
			"mode":      1,
			"timestamp": time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":        1235,
		},
		"motion_detection": map[string]interface{}{
			"status":      true,
			"sensitivity": 80,
			"timestamp":   time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC).UnixMilli(),
			"id":          1236,
		},
	})

	confirmedTime := time.Date(2023, time.October, 25, 15, 0, 0, 0, time.UTC)

	return models.DesiredState{
		ID:            "dev-002", // DeviceID as primary key
		State:         desiredStateJSON,
		Version:       2,
		Status:        true,
		LastDesiredID: 1002,
		ConfirmedAt:   &confirmedTime,
		CreatedAt:     time.Date(2023, time.October, 25, 14, 45, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2023, time.October, 25, 15, 0, 0, 0, time.UTC),
	}
}

func TestGetDesiredState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"desired_states\" WHERE id = \\$1 ORDER BY \"desired_states\".\"id\" LIMIT \\$2").
		WithArgs("dev-001", 1).
		WillReturnRows(mockDesiredStateRows())

	models.MockSetup(db)
	desiredState, Error := testDesiredStateModels.Get("dev-001")
	if Error != nil {
		t.Errorf("Expected a non-nil desired state %s", Error)
	}
	expect := mockDesiredState0()
	assert.Equal(t, expect, desiredState, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetDesiredStateByDeviceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery("SELECT \\* FROM \"desired_states\" WHERE id = \\$1 ORDER BY \"desired_states\".\"id\" LIMIT \\$2").
		WithArgs("dev-001", 1).
		WillReturnRows(mockDesiredStateRows())

	models.MockSetup(db)

	desiredState, Error := testDesiredStateModels.Get("dev-001")
	if Error != nil {
		t.Errorf("Expected a non-nil desired state %s", Error)
	}
	expect := mockDesiredState0()
	assert.Equal(t, expect, desiredState, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDesiredStatesCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	deviceIDs := []string{"dev-001", "dev-002"}
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"desired_states\" WHERE id IN \\(\\$1,\\$2\\)").WithArgs("dev-001", "dev-002").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	models.MockSetup(db)
	count, Error := testDesiredStateModels.Count(deviceIDs)
	if Error != nil {
		t.Errorf("Expected a non-nil count %s", Error)
	}
	expect := int64(2)
	assert.Equal(t, expect, count, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDesiredStatesList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	deviceIDs := []string{"dev-001", "dev-002"}
	mock.ExpectQuery("SELECT \\* FROM \"desired_states\" WHERE id IN \\(\\$1,\\$2\\) ORDER BY created_at DESC LIMIT \\$3").
		WithArgs("dev-001", "dev-002", 10).
		WillReturnRows(mockDesiredStatesRows())

	models.MockSetup(db)
	desiredStates, Error := testDesiredStateModels.List(deviceIDs, 1, 10, nil)
	if Error != nil {
		t.Errorf("Expected a non-nil desired states %s", Error)
	}
	expect := []models.DesiredState{mockDesiredState0(), mockDesiredState1()}
	assert.Equal(t, expect, desiredStates, "they should be equal")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateDesiredState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"desired_states\"").
		WithArgs("dev-001", sqlmock.AnyArg(), 1, true, 1001, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	desiredState := mockDesiredState0()
	err = testDesiredStateModels.Create(&desiredState)
	if err != nil {
		t.Errorf("Expected a non-nil desired state %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateDesiredState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"desired_states\" SET").
		WithArgs(sqlmock.AnyArg(), 1, true, 1001, nil, sqlmock.AnyArg(), sqlmock.AnyArg(), "dev-001").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	desiredState := mockDesiredState0()
	err = testDesiredStateModels.Update(&desiredState)
	if err != nil {
		t.Errorf("Expected a non-nil desired state %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateWithVersionDesiredState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"desired_states\" SET").
		WithArgs("dev-001", sqlmock.AnyArg(), 1, true, 1001, sqlmock.AnyArg(), sqlmock.AnyArg(), "dev-001", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	desiredState := mockDesiredState0()
	version := desiredState.Version
	err = testDesiredStateModels.UpdateWithVersion(&desiredState, version)
	if err != nil {
		t.Errorf("Expected no error when confirming desired state %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteDesiredState(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"desired_states\" WHERE id = \\$1").
		WithArgs("dev-001").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	err = testDesiredStateModels.Delete("dev-001")
	if err != nil {
		t.Errorf("Expected no error when deleting desired state %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
