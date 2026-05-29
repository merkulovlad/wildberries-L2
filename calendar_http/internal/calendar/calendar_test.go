package calendar

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(DateLayout, s)
	if err != nil {
		t.Fatalf("parse %q: %v", s, err)
	}
	return d
}

func TestCreateAndForDay(t *testing.T) {
	c := New()
	day := mustDate(t, "2026-05-29")

	e, err := c.Create(1, day, "standup")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if e.ID == 0 || e.UserID != 1 || e.Text != "standup" {
		t.Fatalf("unexpected event: %+v", e)
	}

	got := c.ForDay(1, day)
	if len(got) != 1 || got[0].ID != e.ID {
		t.Fatalf("ForDay: %+v", got)
	}

	// other user must not see it
	if other := c.ForDay(2, day); len(other) != 0 {
		t.Fatalf("isolation broken: %+v", other)
	}
}

func TestCreateValidation(t *testing.T) {
	c := New()
	day := mustDate(t, "2026-05-29")

	if _, err := c.Create(0, day, "x"); !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("want ErrInvalidUserID, got %v", err)
	}
	if _, err := c.Create(1, day, ""); !errors.Is(err, ErrEmptyEvent) {
		t.Fatalf("want ErrEmptyEvent, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	c := New()
	day := mustDate(t, "2026-05-29")
	e, _ := c.Create(1, day, "old")

	next := mustDate(t, "2026-05-30")
	upd, err := c.Update(e.ID, 1, next, "new")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if upd.Text != "new" || !upd.Date.Equal(next) {
		t.Fatalf("update result: %+v", upd)
	}

	// wrong owner
	if _, err := c.Update(e.ID, 2, next, "x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound for wrong owner, got %v", err)
	}
	// missing id
	if _, err := c.Update(9999, 1, next, "x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound for missing id, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	c := New()
	day := mustDate(t, "2026-05-29")
	e, _ := c.Create(1, day, "x")

	if err := c.Delete(e.ID, 2); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound for wrong owner, got %v", err)
	}
	if err := c.Delete(e.ID, 1); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := c.Delete(e.ID, 1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("double-delete should fail, got %v", err)
	}
}

func TestForWeekAndMonth(t *testing.T) {
	c := New()
	d0 := mustDate(t, "2026-05-01")

	// events at day 0, +3, +10, +40
	for _, off := range []int{0, 3, 10, 40} {
		if _, err := c.Create(1, d0.AddDate(0, 0, off), "e"); err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	if got := c.ForDay(1, d0); len(got) != 1 {
		t.Fatalf("ForDay len=%d", len(got))
	}
	if got := c.ForWeek(1, d0); len(got) != 2 {
		t.Fatalf("ForWeek len=%d, want 2", len(got))
	}
	if got := c.ForMonth(1, d0); len(got) != 3 {
		t.Fatalf("ForMonth len=%d, want 3", len(got))
	}
}

func TestConcurrentCreate(t *testing.T) {
	c := New()
	day := mustDate(t, "2026-05-29")

	const n = 200
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			if _, err := c.Create(1, day, "e"); err != nil {
				t.Errorf("create: %v", err)
			}
		}()
	}
	wg.Wait()

	if got := c.ForDay(1, day); len(got) != n {
		t.Fatalf("got %d events, want %d", len(got), n)
	}
}
