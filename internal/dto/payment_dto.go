package dto

type PaymentDTO struct {
	BookingID uint `json:"booking_id" binding:"required"`
}
