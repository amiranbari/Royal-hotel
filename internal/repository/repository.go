package repository

import (
	"github.com/amiranbari/Royal-hotel/internal/models"
	"time"
)

type DatabaseRepo interface {
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	Authenticate(email, password string, accessLevel int) (int, string, error)
}
