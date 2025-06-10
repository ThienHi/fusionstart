package models

import "time"

type Event struct {
	BaseModel
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Datetime         time.Time `json:"datetime"`
	TotalTickets     uint      `json:"total_tickets" gorm:"not null;check:total_tickets > 0"`
	TicketPrice      float32   `json:"ticket_price" gorm:"not null;check:ticket_price >= 0"`
	AvailableTickets uint      `json:"available_tickets"`
}
