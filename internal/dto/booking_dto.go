package dto

type BookingCreateDTO struct {
	EventId  uint `json:"event_id" binding:"required"`
	UserId   uint `json:"user_id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required"`
}
