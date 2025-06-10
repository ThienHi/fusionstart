package dto

import "time"

type EventCreateDTO struct {
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Datetime     time.Time `json:"datetime" binding:"required"`
	TotalTickets uint      `json:"total_tickets" binding:"required"`
	TicketPrice  float32   `json:"ticket_price" binding:"required"`
}

type EventUpdateDTO struct {
	Name         *string    `json:"name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Datetime     *time.Time `json:"datetime,omitempty"`
	TotalTickets *uint      `json:"total_tickets,omitempty"`
	TicketPrice  *float32   `json:"ticket_price,omitempty"`
}
