package repositories

import (
	"errors"

	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetAll() ([]models.Event, error)
	FindById(id uint) (*models.Event, error)
	Create(event dto.EventCreateDTO) error
	Update(id uint, event dto.EventUpdateDTO) (*models.Event, error)
	Delete(id uint) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	if err := r.db.Find(&events).Order("datetime DESC").Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) FindById(id uint) (*models.Event, error) {
	var event *models.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return event, nil
}

func (r *eventRepository) Create(event dto.EventCreateDTO) error {
	newEvent := models.Event{
		Name:             event.Name,
		Description:      event.Description,
		Datetime:         event.Datetime,
		TotalTickets:     event.TotalTickets,
		TicketPrice:      event.TicketPrice,
		AvailableTickets: event.TotalTickets,
	}
	if err := r.db.Create(&newEvent).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) Update(id uint, data dto.EventUpdateDTO) (*models.Event, error) {
	var event models.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}

	if data.Name != nil {
		event.Name = *data.Name
	}
	if data.Description != nil {
		event.Description = *data.Description
	}
	if data.Datetime != nil {
		event.Datetime = *data.Datetime
	}
	if data.TicketPrice != nil {
		event.TicketPrice = *data.TicketPrice
	}

	if data.TotalTickets != nil {
		var ticketsBooked int
		err := r.db.Model(&models.Booking{}).Where(
			"event_id = ? AND status IN ?",
			id, []string{"PENDING", "CONFIRMED"}).Select("COALESCE(Sum(Quantity), 0)").Scan(&ticketsBooked).Error
		if err != nil {
			return nil, err
		}
		if int(*data.TotalTickets) < ticketsBooked {
			return nil, errors.New("not enough tickets available")
		}

		event.TotalTickets = *data.TotalTickets
		event.AvailableTickets = *data.TotalTickets - uint(ticketsBooked)
	}
	if err := r.db.Save(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
	// return r.db.Model(&models.Event{}).Find(id).Update("is_deleted", true).Error
}
