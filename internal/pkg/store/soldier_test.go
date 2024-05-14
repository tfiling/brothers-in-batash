package store_test

import (
	"brothers_in_batash/internal/pkg/models"
	"brothers_in_batash/internal/pkg/store"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSoldierStore(t *testing.T) {
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)
	assert.NotNil(t, soldierStore)
}

func TestInMemSoldierStore_CreateNewSoldier__success(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldier := models.Soldier{
		ID:             "123",
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	//Act
	err = soldierStore.CreateNewSoldier(soldier)

	//Assert
	assert.NoError(t, err)
}

func TestInMemSoldierStore_CreateNewSoldier__duplicate_id(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldier := models.Soldier{
		ID:             "123",
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	err = soldierStore.CreateNewSoldier(soldier)
	require.NoError(t, err)

	//Act
	err = soldierStore.CreateNewSoldier(soldier)

	//Assert
	assert.Error(t, err)
}

func TestInMemSoldierStore_CreateNewSoldier__invalid_soldier(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	invalidSoldier := models.Soldier{
		ID:             "",
		FirstName:      "",
		LastName:       "",
		PersonalNumber: "",
		Position:       -1,
	}

	//Act
	err = soldierStore.CreateNewSoldier(invalidSoldier)

	//Assert
	assert.Error(t, err)
}

func TestInMemSoldierStore_FindSoldierByID__success(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldier := models.Soldier{
		ID:             "123",
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	err = soldierStore.CreateNewSoldier(soldier)
	require.NoError(t, err)

	//Act
	foundSoldier, err := soldierStore.FindSoldierByID("123")

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, soldier, foundSoldier)
}

func TestInMemSoldierStore_FindSoldierByID__not_found(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	//Act
	_, err = soldierStore.FindSoldierByID("456")

	//Assert
	assert.Error(t, err)
}

func TestInMemSoldierStore_FindAllSoldiers__success(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldiers := []models.Soldier{
		{ID: "1", FirstName: "John", LastName: "Doe", PersonalNumber: "1234567", Position: models.RegularSoldierPosition},
		{ID: "2", FirstName: "Jane", LastName: "Smith", PersonalNumber: "7654321", Position: models.SquadCommanderPosition},
	}

	for _, soldier := range soldiers {
		err := soldierStore.CreateNewSoldier(soldier)
		require.NoError(t, err)
	}

	//Act
	allSoldiers, err := soldierStore.FindAllSoldiers()

	//Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, soldiers, allSoldiers)
}

func TestInMemSoldierStore_UpdateSoldier__success(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldierID := "123"
	soldier := models.Soldier{
		ID:             soldierID,
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	err = soldierStore.CreateNewSoldier(soldier)
	require.NoError(t, err)

	updatedSoldier := soldier
	updatedSoldier.Position = models.SquadCommanderPosition

	//Act
	err = soldierStore.UpdateSoldier(updatedSoldier)

	//Assert
	assert.NoError(t, err)
	foundSoldier, err := soldierStore.FindSoldierByID(soldierID)
	assert.NoError(t, err)
	assert.Equal(t, updatedSoldier, foundSoldier)
}

func TestInMemSoldierStore_UpdateSoldier__not_found(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	//Act
	nonExistentSoldier := models.Soldier{ID: "456"}

	//Assert
	err = soldierStore.UpdateSoldier(nonExistentSoldier)
	assert.Error(t, err)
}

func TestInMemSoldierStore_UpdateSoldier__invalid_data(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldier := models.Soldier{
		ID:             "123",
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	err = soldierStore.CreateNewSoldier(soldier)
	require.NoError(t, err)

	invalidSoldier := soldier
	invalidSoldier.FirstName = ""

	//Act
	err = soldierStore.UpdateSoldier(invalidSoldier)

	//Assert
	assert.Error(t, err)
}

func TestInMemSoldierStore_DeleteSoldier__success(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	soldierID := "123"
	soldier := models.Soldier{
		ID:             soldierID,
		FirstName:      "John",
		LastName:       "Doe",
		PersonalNumber: "1234567",
		Position:       models.RegularSoldierPosition,
	}

	err = soldierStore.CreateNewSoldier(soldier)
	require.NoError(t, err)

	//Act
	err = soldierStore.DeleteSoldier(soldierID)

	//Assert
	assert.NoError(t, err)
	_, err = soldierStore.FindSoldierByID(soldierID)
	assert.Error(t, err)
}

func TestInMemSoldierStore_DeleteSoldier__not_found(t *testing.T) {
	//Arrange
	soldierStore, err := store.NewSoldierStore()
	require.NoError(t, err)

	//Act
	err = soldierStore.DeleteSoldier("456")

	//Assert
	assert.Error(t, err)
}
