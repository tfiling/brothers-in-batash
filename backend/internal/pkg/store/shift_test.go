package store_test

import (
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testSoldier = models.Soldier{
		ID:             "456",
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.SquadCommanderPosition,
		Roles: []models.SoldierRole{{
			ID:   "1",
			Name: "Commander",
		}},
	}
)

func TestNewShiftStore(t *testing.T) {
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)
	assert.NotNil(t, shiftStore)
}

func TestInMemShiftStore_CreateNewShift__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shift := models.Shift{
		ID:        "123",
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		Commander: testSoldier,
	}

	// Act
	err = shiftStore.CreateNewShift(shift)

	// Assert
	assert.NoError(t, err)
}

func TestInMemShiftStore_CreateNewShift__duplicate_id(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shift := models.Shift{
		ID:        "123",
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Commander: testSoldier,
	}

	err = shiftStore.CreateNewShift(shift)
	require.NoError(t, err)

	// Act
	err = shiftStore.CreateNewShift(shift)

	// Assert
	assert.Error(t, err)
}

func TestInMemShiftStore_CreateNewShift__invalid_shift(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	invalidShift := models.Shift{
		ID:   "",
		Name: "",
		Type: -1,
	}

	// Act
	err = shiftStore.CreateNewShift(invalidShift)

	// Assert
	assert.Error(t, err)
}

func TestInMemShiftStore_FindShiftByID__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shift := models.Shift{
		ID:        "123",
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Commander: testSoldier,
	}

	err = shiftStore.CreateNewShift(shift)
	require.NoError(t, err)

	// Act
	foundShifts, err := shiftStore.FindShiftByID("123")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, foundShifts, 1)
	assert.Equal(t, shift, foundShifts[0])
}

func TestInMemShiftStore_FindShiftByID__not_found(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	// Act
	foundShifts, err := shiftStore.FindShiftByID("456")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, foundShifts, 0)
}

func TestInMemShiftStore_FindAllShifts__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shifts := []models.Shift{
		{
			ID:        "1",
			Name:      "Shift 1",
			Type:      models.MotorizedPatrolShiftType,
			StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
			Commander: testSoldier,
		},
		{
			ID:        "2",
			Name:      "Shift 2",
			Type:      models.StaticPostShiftType,
			StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
			Commander: testSoldier,
		},
	}

	for _, shift := range shifts {
		err := shiftStore.CreateNewShift(shift)
		require.NoError(t, err)
	}

	// Act
	allShifts, err := shiftStore.FindAllShifts()

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, shifts, allShifts)
}

func TestInMemShiftStore_UpdateShift__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shiftID := "123"
	shift := models.Shift{
		ID:        shiftID,
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Commander: testSoldier,
	}

	err = shiftStore.CreateNewShift(shift)
	require.NoError(t, err)

	updatedShift := shift
	updatedShift.Name = "Updated Shift"

	// Act
	err = shiftStore.UpdateShift(updatedShift)

	// Assert
	assert.NoError(t, err)
	foundShifts, err := shiftStore.FindShiftByID(shiftID)
	assert.NoError(t, err)
	assert.Len(t, foundShifts, 1)
	assert.Equal(t, updatedShift, foundShifts[0])
}

func TestInMemShiftStore_UpdateShift__not_found(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	// Act
	nonExistentShift := models.Shift{ID: "456"}

	// Assert
	err = shiftStore.UpdateShift(nonExistentShift)
	assert.Error(t, err)
}

func TestInMemShiftStore_UpdateShift__invalid_data(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shift := models.Shift{
		ID:        "123",
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Commander: testSoldier,
	}

	err = shiftStore.CreateNewShift(shift)
	require.NoError(t, err)

	invalidShift := shift
	invalidShift.Name = ""

	// Act
	err = shiftStore.UpdateShift(invalidShift)

	// Assert
	assert.Error(t, err)
}

func TestInMemShiftStore_DeleteShift__success(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	shiftID := "123"
	shift := models.Shift{
		ID:        shiftID,
		Name:      "Test Shift",
		Type:      models.MotorizedPatrolShiftType,
		StartTime: time.Date(2025, time.April, 9, 15, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 9, 16, 0, 0, 0, time.UTC),
		Commander: testSoldier,
	}

	err = shiftStore.CreateNewShift(shift)
	require.NoError(t, err)

	// Act
	err = shiftStore.DeleteShift(shiftID)

	// Assert
	assert.NoError(t, err)
	foundShifts, err := shiftStore.FindShiftByID(shiftID)
	assert.NoError(t, err)
	assert.Len(t, foundShifts, 0)
}

func TestInMemShiftStore_DeleteShift__not_found(t *testing.T) {
	// Arrange
	shiftStore, err := store.NewShiftStore()
	require.NoError(t, err)

	// Act
	err = shiftStore.DeleteShift("456")

	// Assert
	assert.Error(t, err)
}
