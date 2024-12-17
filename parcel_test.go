package main

import (
	"math/rand"
	"testing"
	"time"

	"database/sql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	parcel.Number = id
	parcelAdded, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel, parcelAdded)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	assert.NotEmpty(t, err)
	assert.Equal(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	parcel.Number = id

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	parcelNewAddress, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcelNewAddress.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	parcel.Number = id

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	newParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newParcel.Status, ParcelStatusSent)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)

	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, storedParcels, len(parcelMap))

	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		assert.True(t, ok)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
