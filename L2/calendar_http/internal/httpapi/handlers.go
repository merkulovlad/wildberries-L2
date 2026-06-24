package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/merkulovlad/wildberries-L2/calendar_http/internal/calendar"
)

// Server wires HTTP routes to the calendar service.
type Server struct {
	cal *calendar.Calendar
	mux *http.ServeMux
}

func NewServer(cal *calendar.Calendar) *Server {
	s := &Server{cal: cal, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("/create_event", s.methodOnly(http.MethodPost, s.handleCreate))
	s.mux.HandleFunc("/update_event", s.methodOnly(http.MethodPost, s.handleUpdate))
	s.mux.HandleFunc("/delete_event", s.methodOnly(http.MethodPost, s.handleDelete))
	s.mux.HandleFunc("/events_for_day", s.methodOnly(http.MethodGet, s.handleForDay))
	s.mux.HandleFunc("/events_for_week", s.methodOnly(http.MethodGet, s.handleForWeek))
	s.mux.HandleFunc("/events_for_month", s.methodOnly(http.MethodGet, s.handleForMonth))
}

func (s *Server) methodOnly(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h(w, r)
	}
}

// eventInput is the parsed body of create/update/delete requests.
type eventInput struct {
	ID     int64
	UserID int64
	Date   time.Time
	Text   string
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	in, err := parseBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	e, err := s.cal.Create(in.UserID, in.Date, in.Text)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeResult(w, e)
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	in, err := parseBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.ID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	e, err := s.cal.Update(in.ID, in.UserID, in.Date, in.Text)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeResult(w, e)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	in, err := parseBody(r)
	// Delete only needs id+user_id; tolerate missing date/text.
	if err != nil && !errors.Is(err, errMissingDate) && !errors.Is(err, errMissingText) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.ID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if in.UserID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	if err := s.cal.Delete(in.ID, in.UserID); err != nil {
		writeBusinessError(w, err)
		return
	}
	writeResult(w, "deleted")
}

func (s *Server) handleForDay(w http.ResponseWriter, r *http.Request) {
	s.handleQuery(w, r, s.cal.ForDay)
}

func (s *Server) handleForWeek(w http.ResponseWriter, r *http.Request) {
	s.handleQuery(w, r, s.cal.ForWeek)
}

func (s *Server) handleForMonth(w http.ResponseWriter, r *http.Request) {
	s.handleQuery(w, r, s.cal.ForMonth)
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request,
	fn func(int64, time.Time) []calendar.Event) {
	q := r.URL.Query()
	userID, err := parseInt(q.Get("user_id"))
	if err != nil || userID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}
	date, err := time.Parse(calendar.DateLayout, q.Get("date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date")
		return
	}
	writeResult(w, fn(userID, date))
}

// --- parsing ---

var (
	errMissingDate = errors.New("invalid date")
	errMissingText = errors.New("event text must not be empty")
)

// parseBody reads either application/json or url-encoded form into eventInput.
// Returned errors are safe to surface verbatim to clients.
func parseBody(r *http.Request) (eventInput, error) {
	var in eventInput
	ct := r.Header.Get("Content-Type")

	if strings.HasPrefix(ct, "application/json") {
		var body struct {
			ID     int64  `json:"id"`
			UserID int64  `json:"user_id"`
			Date   string `json:"date"`
			Text   string `json:"event"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return in, errors.New("invalid json")
		}
		in.ID = body.ID
		in.UserID = body.UserID
		in.Text = body.Text
		if body.UserID <= 0 {
			return in, errors.New("invalid user_id")
		}
		if body.Date == "" {
			return in, errMissingDate
		}
		d, err := time.Parse(calendar.DateLayout, body.Date)
		if err != nil {
			return in, errMissingDate
		}
		in.Date = d
		if body.Text == "" {
			return in, errMissingText
		}
		return in, nil
	}

	if err := r.ParseForm(); err != nil {
		return in, errors.New("invalid form")
	}
	if v := r.PostFormValue("id"); v != "" {
		id, err := parseInt(v)
		if err != nil {
			return in, errors.New("invalid id")
		}
		in.ID = id
	}
	uid, err := parseInt(r.PostFormValue("user_id"))
	if err != nil || uid <= 0 {
		return in, errors.New("invalid user_id")
	}
	in.UserID = uid

	dateStr := r.PostFormValue("date")
	if dateStr == "" {
		return in, errMissingDate
	}
	d, err := time.Parse(calendar.DateLayout, dateStr)
	if err != nil {
		return in, errMissingDate
	}
	in.Date = d

	in.Text = r.PostFormValue("event")
	if in.Text == "" {
		return in, errMissingText
	}
	return in, nil
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// --- response helpers ---

type okResp struct {
	Result any `json:"result"`
}

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeResult(w http.ResponseWriter, v any) {
	writeJSON(w, http.StatusOK, okResp{Result: v})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errResp{Error: msg})
}

// writeBusinessError maps domain errors to HTTP status codes per spec:
// 503 for business-logic errors, 500 for anything unexpected.
func writeBusinessError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calendar.ErrNotFound),
		errors.Is(err, calendar.ErrInvalidUserID),
		errors.Is(err, calendar.ErrEmptyEvent),
		errors.Is(err, calendar.ErrInvalidDate):
		writeError(w, http.StatusServiceUnavailable, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
