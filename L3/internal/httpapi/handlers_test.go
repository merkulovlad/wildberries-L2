package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
)

func TestCreateNotification(t *testing.T) {
	sendAt := time.Date(2026, 6, 24, 15, 0, 0, 0, time.UTC)
	service := fakeService{
		createFn: func(_ context.Context, in notification.CreateNotification) (notification.Notification, error) {
			if in.Recipient != "user-1" {
				t.Fatalf("recipient = %q, want user-1", in.Recipient)
			}
			if in.Message != "hello" {
				t.Fatalf("message = %q, want hello", in.Message)
			}
			if !in.SendAt.Equal(sendAt) {
				t.Fatalf("send_at = %s, want %s", in.SendAt, sendAt)
			}
			return notification.Notification{ID: "n-1", Status: notification.StatusScheduled}, nil
		},
	}

	body := bytes.NewBufferString(`{"recipient":"user-1","message":"hello","send_at":"2026-06-24T15:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/notify", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	NewRouter(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var got responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Result.ID != "n-1" || got.Result.Status != "scheduled" {
		t.Fatalf("result = %+v, want created status", got.Result)
	}
}

func TestCreateNotificationValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(`{"message":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	NewRouter(fakeService{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetNotification(t *testing.T) {
	service := fakeService{
		getFn: func(_ context.Context, id string) (notification.Notification, error) {
			if id != "n-1" {
				t.Fatalf("id = %q, want n-1", id)
			}
			return notification.Notification{ID: id, Status: notification.StatusScheduled}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/notify/n-1", nil)
	rec := httptest.NewRecorder()

	NewRouter(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestCancelNotification(t *testing.T) {
	service := fakeService{
		cancelFn: func(_ context.Context, id string) (notification.Notification, error) {
			if id != "n-1" {
				t.Fatalf("id = %q, want n-1", id)
			}
			return notification.Notification{ID: id, Status: notification.StatusCanceled}, nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/notify/n-1", nil)
	rec := httptest.NewRecorder()

	NewRouter(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

type fakeService struct {
	createFn func(context.Context, notification.CreateNotification) (notification.Notification, error)
	getFn    func(context.Context, string) (notification.Notification, error)
	cancelFn func(context.Context, string) (notification.Notification, error)
}

func (s fakeService) Create(ctx context.Context, in notification.CreateNotification) (notification.Notification, error) {
	if s.createFn == nil {
		return notification.Notification{}, notification.ErrNotImplemented
	}
	return s.createFn(ctx, in)
}

func (s fakeService) Get(ctx context.Context, id string) (notification.Notification, error) {
	if s.getFn == nil {
		return notification.Notification{}, notification.ErrNotImplemented
	}
	return s.getFn(ctx, id)
}

func (s fakeService) Cancel(ctx context.Context, id string) (notification.Notification, error) {
	if s.cancelFn == nil {
		return notification.Notification{}, notification.ErrNotImplemented
	}
	return s.cancelFn(ctx, id)
}

type responseEnvelope struct {
	Result notificationResponse `json:"result"`
}
