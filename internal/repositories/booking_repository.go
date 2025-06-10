package repositories

import (
	"errors"
	"time"

	"github.com/thienhi/fusionstart/internal/constants"
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	GetAll() ([]models.Booking, error)
	Create(booking dto.BookingCreateDTO) (*models.Booking, error)
	Delete(id uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) GetAll() ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) Create(bookingInput dto.BookingCreateDTO) (*models.Booking, error) {
	var (
		event         = models.Event{}
		user          = models.User{}
		ticketsBooked uint
		booking       *models.Booking
	)

	if err := r.db.First(&user, bookingInput.UserId).Error; err != nil {
		return nil, err
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&event, "id = ?", bookingInput.EventId).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Booking{}).Where(
			"event_id = ? AND status IN ?",
			bookingInput.EventId,
			[]string{"PENDING", "CONFIRMED"}).Select("COALESCE(Sum(Quantity), 0)").Scan(&ticketsBooked).Error; err != nil {
			return err
		}

		available := int(event.TotalTickets) - int(ticketsBooked)
		if available < int(bookingInput.Quantity) {
			return errors.New("not enough tickets available")
		}

		booking = &models.Booking{
			UserID:      bookingInput.UserId,
			EventID:     bookingInput.EventId,
			Quantity:    bookingInput.Quantity,
			Status:      constants.BookingStatusPending,
			TotalAmount: float32(bookingInput.Quantity) * event.TicketPrice,
			ExpiresAt: func() time.Time {
				now := time.Now()
				return now.Add(time.Minute * 15)
			}(),
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		event.AvailableTickets -= booking.Quantity
		if err := tx.Save(&event).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, "id = ?", id).Error
}
