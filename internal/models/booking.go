package models

import (
	"time"

	"github.com/thienhi/fusionstart/internal/constants"
)

type Booking struct {
	BaseModel
	// Relations to User
	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// Relations to Event
	EventID uint  `json:"event_id" gorm:"not null;index"`
	Event   Event `json:"event" gorm:"foreignKey:EventID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// For Booking
	Status      constants.BookingStatus `json:"status" gorm:"type:varchar(20);default:'PENDING';index"`
	Quantity    uint                    `json:"quantity"`
	TotalAmount float32                 `json:"total_amount"`
	ExpiresAt   time.Time               `json:"expires_at" gorm:"not null;index"`
	ConfirmedAt *time.Time              `json:"confirmed_at"`
	CancelledAt *time.Time              `json:"cancelled_at"`
}
