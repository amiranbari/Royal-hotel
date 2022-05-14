package handlers

import (
	"errors"
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
		m.App.Session.Put(r.Context(), "error", "start date and end date required!")
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

func (m *Repository) Login(rw http.ResponseWriter, r *http.Request) {
	renders.Template(rw, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(rw http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "password")
	form.IsEmail("username")
	form.MinLength("password", 8)
	if !form.Valid() {
		renders.Template(rw, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	email := r.Form.Get("username")
	password := r.Form.Get("password")

	id, _, err := m.DB.Authenticate(email, password, 5)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	if r.Form.Get("redirect") != "" {
		http.Redirect(rw, r, r.Form.Get("redirect"), http.StatusSeeOther)
		return
	}

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

// Logout users
func (m *Repository) Logout(rw http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (m *Repository) ChooseRoom(rw http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
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

	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room = room
	data["reservation"] = res
	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	renders.Template(rw, r, "single-room.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

func (m *Repository) MakeReservation(rw http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data["reservation"] = res
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	renders.Template(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(rw, errors.New("can't parse form"))
		return
	}

	form := forms.New(r.PostForm)
	form.Required("firstname", "lastname", "email", "phone")
	form.IsEmail("email")

	if !form.Valid() {
		helpers.ServerError(rw, errors.New("form is not valid"))
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(rw, errors.New("there is a problem in your reservation"))
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("firstname"),
		LastName:  r.Form.Get("lastname"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		RoomId:    res.Room.ID,
	}

	newReservationId, err := m.DB.InsertReservation(reservation)

	if err != nil {
		helpers.ServerError(rw, err)
		return
	}

	restriction := models.RoomRestriction{
		RoomId:        res.Room.ID,
		ReservationId: newReservationId,
		RestrictionId: 1,
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
	}

	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		helpers.ServerError(rw, errors.New("can't insert restriction to database"))
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	http.Redirect(rw, r, "/profile", http.StatusSeeOther)

}
