package dbrepo

import (
	"github.com/amiranbari/Royal-hotel/internal/models"
	"time"
)

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	return room, nil
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	var newId int
	return newId, nil
}

func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, password string, accessLevel int) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}
