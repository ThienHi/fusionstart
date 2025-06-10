package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thienhi/fusionstart/internal/constants"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/models"
	"github.com/thienhi/fusionstart/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&models.BaseModel{},
		&models.Booking{},
		&models.User{},
		&models.Event{},
	)
	if err != nil {
		t.Fatalf("Database migration error: %v", err)
	}

	return db
}

func setupUserAndEvent(db *gorm.DB) (models.User, models.Event) {
	user := models.User{
		Name:     "Name A",
		Email:    "a00@example.com",
		Password: "password",
	}
	db.Create(&user)

	event := models.Event{
		Name:             "Concert A",
		TotalTickets:     uint(100),
		TicketPrice:      50.00,
		Description:      "A great concert",
		Datetime:         time.Now().Add(24 * time.Hour),
		AvailableTickets: uint(100),
	}
	db.Create(&event)

	return user, event
}

func TestCreateBookingFail(t *testing.T) {
	db := setupTestDB(t)
	bookingRepositories := repositories.NewBookingRepository(db)
	_, err := bookingRepositories.Create(dto.BookingCreateDTO{})
	assert.Error(t, err)
}

func TestCreateBookingSuccess(t *testing.T) {
	db := setupTestDB(t)

	user, event := setupUserAndEvent(db)

	bookingRepositories := repositories.NewBookingRepository(db)

	b, err := bookingRepositories.Create(dto.BookingCreateDTO{
		UserId:   user.ID,
		EventId:  event.ID,
		Quantity: 2,
	})

	assert.NoError(t, err)
	assert.Equal(t, user.ID, b.UserID)
	assert.Equal(t, uint(2), b.Quantity)
	assert.Equal(t, constants.BookingStatusPending, b.Status)
	assert.NotEqual(t, constants.BookingStatusCancelled, b.Status)
}

func TestCreateBookingFailNotEnoughTickets(t *testing.T) {
	db := setupTestDB(t)
	user, event := setupUserAndEvent(db)

	fmt.Println("event ------ ", event)
	fmt.Println("user ------ ", user)

	// First booking - 90 tickets
	bookingRepositories := repositories.NewBookingRepository(db)
	_, err1 := bookingRepositories.Create(dto.BookingCreateDTO{
		UserId:   user.ID,
		EventId:  event.ID,
		Quantity: uint(90),
	})
	assert.NoError(t, err1)

	// Second booking - 30 tickets (should fail, only 10 tickets left)
	_, err2 := bookingRepositories.Create(dto.BookingCreateDTO{
		UserId:   user.ID,
		EventId:  event.ID,
		Quantity: uint(30),
	})

	assert.Error(t, err2)
	assert.EqualError(t, err2, "not enough tickets available")
}
