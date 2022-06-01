package dbrepo

import (
	"errors"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"time"
)

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	if end.Format("2006-01-02") == "2050-01-03" {
		return rooms, errors.New("Some error ...")
	}
	if end.Format("2006-01-02") == "2050-01-05" {
		rooms = []models.Room{
			{
				1,
				"test-room",
				time.Now(),
				time.Now(),
				3500,
			},
		}
	}
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
