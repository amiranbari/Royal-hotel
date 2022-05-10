package dbrepo

import (
	"context"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"time"
)

func (m *PostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
				select r.id, r.title, r.price
			from 
				rooms r
				where r.id not in 
				(select rr.room_id from room_restrictions rr where $1 <= end_date and $2 >= start_date) 
			`

	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Title, &room.Price)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}
