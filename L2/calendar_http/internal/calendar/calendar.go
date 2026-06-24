// Package calendar holds the in-memory event store and the business-logic
// operations. It does not depend on net/http so it can be tested in isolation.
package calendar

import (
	"errors"
	"sync"
	"time"
)

// DateLayout is the canonical YYYY-MM-DD format accepted on the wire.
const DateLayout = "2006-01-02"

var (
	ErrNotFound      = errors.New("event not found")
	ErrInvalidUserID = errors.New("invalid user_id")
	ErrInvalidDate   = errors.New("invalid date")
	ErrEmptyEvent    = errors.New("event text must not be empty")
)

// Event is one calendar entry. Date is normalized to midnight UTC so that
// per-day grouping is unambiguous.
type Event struct {
	ID     int64     `json:"id"`
	UserID int64     `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"event"`
}

// Calendar is a concurrency-safe in-memory event store keyed by event ID.
type Calendar struct {
	mu     sync.RWMutex
	nextID int64
	events map[int64]Event
}

func New() *Calendar {
	return &Calendar{events: make(map[int64]Event)}
}

// Create stores a new event and returns it with the assigned ID.
func (c *Calendar) Create(userID int64, date time.Time, text string) (Event, error) {
	if userID <= 0 {
		return Event{}, ErrInvalidUserID
	}
	if text == "" {
		return Event{}, ErrEmptyEvent
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.nextID++
	e := Event{
		ID:     c.nextID,
		UserID: userID,
		Date:   normalize(date),
		Text:   text,
	}
	c.events[e.ID] = e
	return e, nil
}

// Update replaces date/text of the event identified by id. The owner must
// match userID to prevent cross-user overwrites.
func (c *Calendar) Update(id, userID int64, date time.Time, text string) (Event, error) {
	if userID <= 0 {
		return Event{}, ErrInvalidUserID
	}
	if text == "" {
		return Event{}, ErrEmptyEvent
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.events[id]
	if !ok || e.UserID != userID {
		return Event{}, ErrNotFound
	}
	e.Date = normalize(date)
	e.Text = text
	c.events[id] = e
	return e, nil
}

// Delete removes the event if it exists and belongs to userID.
func (c *Calendar) Delete(id, userID int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.events[id]
	if !ok || e.UserID != userID {
		return ErrNotFound
	}
	delete(c.events, id)
	return nil
}

// ForDay returns all of userID's events on the given day.
func (c *Calendar) ForDay(userID int64, day time.Time) []Event {
	day = normalize(day)
	return c.inRange(userID, day, day.AddDate(0, 0, 1))
}

// ForWeek returns events in the 7-day window starting on day (inclusive).
func (c *Calendar) ForWeek(userID int64, day time.Time) []Event {
	day = normalize(day)
	return c.inRange(userID, day, day.AddDate(0, 0, 7))
}

// ForMonth returns events in the 30-day window starting on day (inclusive).
// Spec is ambiguous between calendar-month and rolling-30-day; we pick the
// latter so that any starting date works consistently.
func (c *Calendar) ForMonth(userID int64, day time.Time) []Event {
	day = normalize(day)
	return c.inRange(userID, day, day.AddDate(0, 1, 0))
}

// inRange returns events with start <= e.Date < end for userID.
func (c *Calendar) inRange(userID int64, start, end time.Time) []Event {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]Event, 0)
	for _, e := range c.events {
		if e.UserID != userID {
			continue
		}
		if (e.Date.Equal(start) || e.Date.After(start)) && e.Date.Before(end) {
			out = append(out, e)
		}
	}
	return out
}

func normalize(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
