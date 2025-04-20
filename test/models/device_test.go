package models_testing

import (
	"encoding/json"
	"nexus/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var testDeivcesModels models.DeviceModels

func mockDeviceRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "license_id", "name", "product_id", "group_id", "features", "state", "version", "sdk_version", "ip",
		"online", "location", "created_at", "active_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockDevice0())...)
}

func mockDevicesRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "license_id", "name", "product_id", "group_id", "features", "state", "version", "sdk_version", "ip",
		"online", "location", "created_at", "active_at", "updated_at", "status"}).
		AddRow(StructToSlice(mockDevice0())...).
		AddRow(StructToSlice(mockDevice1())...)
}

func mockDevice0() models.Device {
	faetures, _ := json.Marshal([]string{"p2p", "cloud_storage"})
	state, _ := json.Marshal(map[string]interface{}{"volume": 100, "record": map[string]interface{}{"status": true, "mode": 1}})

	return models.Device{
		ID:         "d0e1f2a3-b4c5-6d7e-8f9a-0b1c2d3e4f5a",
		LicenseID:  "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
		Name:       "Alice",
		ProductID:  "1",
		GroupID:    "1010",
		Features:   faetures,
		State:      state,
		Version:    "v0.1.1",
		SDKVersion: "sdk_0.0.1",
		IP:         "192.168.22.123",
		Online:     true,
		Location:   "121.123456,31.123456",
		CreatedAt:  time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		ActiveAt:   time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2020, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:     true,
	}
}

func mockDevice1() models.Device {
	featrues, _ := json.Marshal([]string{"p2p", "cloud_storage"})
	state, _ := json.Marshal(map[string]interface{}{"volume": 100, "record": map[string]interface{}{"status": true, "mode": 1}})

	return models.Device{
		ID:         "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
		LicenseID:  "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d",
		Name:       "Dave",
		ProductID:  "f0e1d2c3-b4a5-6b7c-8d9e-0f1a2b3c4d5e",
		GroupID:    "1010",
		Features:   featrues,
		State:      state,
		Version:    "v0.1.1",
		SDKVersion: "sdk_0.0.1",
		IP:         "192.168.12.13",
		Online:     true,
		Location:   "34.0522,-118.2437",
		CreatedAt:  time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC),
		ActiveAt:   time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2023, time.October, 25, 14, 30, 0, 0, time.UTC),
		Status:     true,
	}
}

func TestGetDevice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM \"devices\" WHERE id = \\$1").
		WithArgs("d0e1f2a3-b4c5-6d7e-8f9a-0b1c2d3e4f5a").
		WillReturnRows(mockDeviceRows())

	models.MockSetup(db)
	device, Error := testDeivcesModels.Get("d0e1f2a3-b4c5-6d7e-8f9a-0b1c2d3e4f5a")
	if Error != nil {
		t.Errorf("Expected a non-nil device %s", Error)
	}
	expect := mockDevice0()
	assert.Equal(t, expect, device, "they should be equal")
}

func TestGetDevices(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT \\* FROM \"devices\" WHERE group_id IN \\(\\$1\\)").
		WithArgs(1010).
		WillReturnRows(mockDevicesRows())

	models.MockSetup(db)
	devices, Error := testDeivcesModels.List([]int32{1010})
	if Error != nil {
		t.Errorf("Expected a non-nil device %s", Error)
	}
	expect := []models.Device{mockDevice0(), mockDevice1()}
	assert.Equal(t, expect, devices, "they should be equal")
}

func TestCreateDevice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"devices\"").
		WithArgs(StructToSlice(mockDevice0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	device := mockDevice0()
	err = testDeivcesModels.Create(&device)
	if err != nil {
		t.Errorf("Expected a non-nil device %s", err)
	}
}

func TestUpdateDevice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE \"devices\" SET").
		WithArgs(StructToUpdateSlice(mockDevice0())...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	models.MockSetup(db)
	device := mockDevice0()
	err = testDeivcesModels.Update(&device)
	if err != nil {
		t.Errorf("Expected a non-nil device %s", err)
	}
}
