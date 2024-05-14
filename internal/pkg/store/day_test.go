package store_test

import (
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCommander = models.Soldier{
	ID:             "1",
	FirstName:      "Gal",
	LastName:       "Tfilin",
	PersonalNumber: "1212121",
	Position:       models.RegularSoldierPosition,
}

var testDaySchedule = models.DaySchedule{
	Date: time.Now(),
	Shifts: []models.Shift{
		{
			ID:   "1",
			Name: "Morning Shift",
			Type: models.DailyDutyShiftType,
			StartTime: models.TimeOfDay{
				Hour:   0,
				Minute: 0,
			},
			EndTime: models.TimeOfDay{
				Hour:   1,
				Minute: 0,
			},
			Commander: testCommander,
		},
	},
}

func TestInMemDaySchedStore_CreateNewDaySchedule__success(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	assert.NoError(t, err)

	// Act
	err = dayStore.CreateNewDaySchedule(testDaySchedule)

	// Assert
	assert.NoError(t, err)
	storedDaySchedule, err := dayStore.FindDaySchedule(testDaySchedule.Date)
	require.NoError(t, err)
	require.Len(t, storedDaySchedule, 1)
	assert.Equal(t, testDaySchedule, storedDaySchedule[0])
}

func TestInMemDaySchedStore_CreateNewDaySchedule__error_on_invalid_data(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	assert.NoError(t, err)
	daySchedule := models.DaySchedule{
		Date: time.Now(),
	}

	// Act
	err = dayStore.CreateNewDaySchedule(daySchedule)

	// Assert
	assert.Error(t, err)
	storedDaySchedule, err := dayStore.FindDaySchedule(daySchedule.Date)
	require.NoError(t, err)
	assert.Empty(t, storedDaySchedule)
}

func TestInMemDaySchedStore_FindDaySchedule__found(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	require.NoError(t, err)
	err = dayStore.CreateNewDaySchedule(testDaySchedule)
	require.NoError(t, err)

	// Act
	result, err := dayStore.FindDaySchedule(testDaySchedule.Date)

	// Assert
	assert.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, testDaySchedule, result[0])
}

func TestInMemDaySchedStore_FindDaySchedule__not_found(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	require.NoError(t, err)
	require.NotNil(t, dayStore)
	date := time.Date(2023, time.May, 15, 0, 0, 0, 0, time.UTC)

	// Act
	result, err := dayStore.FindDaySchedule(date)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestInMemDaySchedStore_UpdateDaySchedule__success(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	require.NoError(t, err)
	require.NotNil(t, dayStore)
	err = dayStore.CreateNewDaySchedule(testDaySchedule)
	require.NoError(t, err)
	updatedDaySchedule := testDaySchedule
	updatedDaySchedule.Shifts = append(updatedDaySchedule.Shifts, updatedDaySchedule.Shifts[0])
	updatedDaySchedule.Shifts[1].StartTime.Hour += 2
	updatedDaySchedule.Shifts[1].EndTime.Hour += 3

	// Act
	err = dayStore.UpdateDaySchedule(updatedDaySchedule)

	// Assert
	assert.NoError(t, err)
	storedDaySchedule, err := dayStore.FindDaySchedule(testDaySchedule.Date)
	require.NoError(t, err)
	require.Len(t, storedDaySchedule, 1)
	assert.Equal(t, updatedDaySchedule, storedDaySchedule[0])
}

func TestInMemDaySchedStore_UpdateDaySchedule__not_found(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	require.NoError(t, err)
	require.NotNil(t, dayStore)

	// Act
	err = dayStore.UpdateDaySchedule(testDaySchedule)

	// Assert
	assert.Error(t, err)
	storedDaySchedule, err := dayStore.FindDaySchedule(testDaySchedule.Date)
	assert.NoError(t, err)
	assert.Empty(t, storedDaySchedule)
}

func TestInMemDaySchedStore_UpdateDaySchedule__invalid_schedule(t *testing.T) {
	// Arrange
	dayStore, err := store.NewInMemDaySchedStore()
	require.NoError(t, err)
	require.NotNil(t, dayStore)
	err = dayStore.CreateNewDaySchedule(testDaySchedule)
	require.NoError(t, err)
	invalidDaySchedule := models.DaySchedule{
		Date: testDaySchedule.Date,
		Shifts: []models.Shift{
			{ID: "", Name: ""}, // Invalid shift
		},
	}

	// Act
	err = dayStore.UpdateDaySchedule(invalidDaySchedule)

	// Assert
	assert.Error(t, err)
	storedDaySchedule, err := dayStore.FindDaySchedule(testDaySchedule.Date)
	require.NoError(t, err)
	require.Len(t, storedDaySchedule, 1)
	assert.Equal(t, testDaySchedule, storedDaySchedule[0]) // Original schedule unchanged
}
