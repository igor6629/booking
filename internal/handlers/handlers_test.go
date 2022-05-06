package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/igor6629/booking/internal/driver"
	"github.com/igor6629/booking/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"vc", "/vega-cottage", "GET", http.StatusOK},
	{"sh", "/summer-house", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	/*{"sap", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"sapj", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"mr", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Vega Cottage",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong responce code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong responce code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong responce code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong responce code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// testing for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong responce code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for invalid start_date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong responce code for invalid start_date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for invalid end_date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong responce code for invalid end_date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for invalid room_id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong responce code for invalid room_id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong responce code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// testing for failure to insert reservation to database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// testing for failure to insert restriction to database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")

	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no rooms available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- rooms are available
	*****************************************/
	// this time, we specify a start date before 2040-01-01, which will give us
	// a non-empty slice, indicating that rooms are available
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")

	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// third case -- empty post body
	*****************************************/
	// create our request with a nil body, so parsing form fails
	req, _ = http.NewRequest("POST", "/search-availability", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// fourth case -- start date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// fifth case -- end date in wrong format
	*****************************************/
	// this time, we specify a start date in the wrong format
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "invalid")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	/*****************************************
	// sixth case -- database query fails
	*****************************************/
	// this time, we specify a start date of 2060-01-01, which will cause
	// our testdb repo to return an error
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	// since we have rooms available, we expect to get status http.StatusTemporaryRedirect
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Errorf("failed to parse json")
	}

	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	/*****************************************
	// second case -- rooms not available
	*****************************************/
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if !j.OK {
		t.Error("Got no availability when some was expected in AvailabilityJSON")
	}

	/*****************************************
	// third case -- no request body
	*****************************************/
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if j.OK || j.Message != "Internal server error" {
		t.Error("Got availability when request body was empty")
	}

	/*****************************************
	// fourth case -- database error
	*****************************************/
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.AvailabilityJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date < 2049-12-31, we expect availability
	if j.OK || j.Message != "Error querying database" {
		t.Error("Got availability when simulating database error")
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Vega Cottage",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	/*****************************************
	// second case -- reservation not in session
	*****************************************/
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Vega Cottage",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// second case -- reservation not in session
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// third case -- missing url parameter, or malformed parameter
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	/*****************************************
	// first case -- database works
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Vega Cottage",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	/*****************************************
	// second case -- database failed
	*****************************************/
	req, _ = http.NewRequest("GET", "/book-room?s=2040-01-01&e=2040-01-02&id=4", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
