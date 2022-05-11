package handlers

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/driver"
	"github.com/amiranbari/Royal-hotel/internal/forms"
	"github.com/amiranbari/Royal-hotel/internal/helpers"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"github.com/amiranbari/Royal-hotel/internal/renders"
	"github.com/amiranbari/Royal-hotel/internal/repository"
	"github.com/amiranbari/Royal-hotel/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Index(rw http.ResponseWriter, r *http.Request) {
	renders.Template(rw, r, "index.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SearchForRoom(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "form is not valid!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.Form)
	form.Required("start_date", "end_date")
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "form is not valid!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:0405PM '06 --700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't search in availability rooms!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No available room!")
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["end_date"] = sd
	stringMap["start_date"] = ed

	renders.Template(rw, r, "search.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(rw, r, "/search", http.StatusTemporaryRedirect)
		return
	}

	res.RoomId = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) MakeReservation(rw http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.Title = room.Title
	data["reservation"] = res

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Put(r.Context(), "reservation", res)

	renders.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(rw, r, "/make-reservation", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("firstname", "lastname", "email", "phone")
	form.IsEmail("email")

	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Form is not valid!")
		http.Redirect(rw, r, "/make-reservation", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("firstname"),
		LastName:  r.Form.Get("lastname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		RoomId:    res.RoomId,
	}
	reservation.Room.Title = res.Room.Title

	newReservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation to database!")
		http.Redirect(rw, r, "/make-reservation", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		RoomId:        res.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert restriction to database!")
		http.Redirect(rw, r, "/make-reservation", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	//send reservation mail

	//html := fmt.Sprintf(`
	//	<strong>Reservation Confirmation</stronge><br>
	//	Dear %s: <br>
	//	This is to confirm your reservation from %s to %s.
	//	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	//
	//msg := models.MailData{
	//	To:      reservation.Email,
	//	From:    "me@here.com",
	//	Subject: "Reservation confirmation",
	//	Content: html,
	//}
	//
	//m.App.MailChan <- msg

	//send room mail
	//html = fmt.Sprintf(`
	//	<strong>Reservation Notification</stronge><br>
	//	A reservation has been made for %s from %s to %s.
	//	`, reservation.Room.Title, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	//
	//msg = models.MailData{
	//	To:      reservation.Email,
	//	From:    "me@here.com",
	//	Subject: "Reservation confirmation",
	//	Content: html,
	//}
	//
	//m.App.MailChan <- msg

	http.Redirect(rw, r, "/reservation", http.StatusSeeOther)

}
