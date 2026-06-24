package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service notification.Service
}

func NewRouter(service notification.Service) http.Handler {
	router := ginext.New("release")
	router.Use(ginext.Logger(), ginext.Recovery())

	handler := &Handler{service: service}
	router.GET("/", handler.index)
	router.POST("/notify", handler.create)
	router.GET("/notify/:id", handler.get)
	router.DELETE("/notify/:id", handler.cancel)

	return router
}

func (h *Handler) index(c *ginext.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
}

func (h *Handler) create(c *ginext.Context) {
	var req createNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid json body")
		return
	}

	in, err := req.toInput()
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	notificationEntity, err := h.service.Create(c.Request.Context(), in)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response{Result: newNotificationResponse(notificationEntity)})
}

func (h *Handler) get(c *ginext.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		writeError(c, http.StatusBadRequest, "id is required")
		return
	}

	notificationEntity, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response{Result: newNotificationResponse(notificationEntity)})
}

func (h *Handler) cancel(c *ginext.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		writeError(c, http.StatusBadRequest, "id is required")
		return
	}

	notificationEntity, err := h.service.Cancel(c.Request.Context(), id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, response{Result: newNotificationResponse(notificationEntity)})
}

type createNotificationRequest struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	SendAt    string `json:"send_at"`
}

func (r createNotificationRequest) toInput() (notification.CreateNotification, error) {
	var in notification.CreateNotification

	in.Recipient = strings.TrimSpace(r.Recipient)
	if in.Recipient == "" {
		return in, errors.New("recipient is required")
	}

	in.Message = strings.TrimSpace(r.Message)
	if in.Message == "" {
		return in, errors.New("message is required")
	}

	if strings.TrimSpace(r.SendAt) == "" {
		return in, errors.New("send_at is required")
	}

	sendAt, err := time.Parse(time.RFC3339, r.SendAt)
	if err != nil {
		return in, errors.New("send_at must be RFC3339")
	}
	in.SendAt = sendAt

	return in, nil
}

type response struct {
	Result any `json:"result"`
}

type notificationResponse struct {
	ID        string     `json:"id"`
	Status    string     `json:"status"`
	Recipient string     `json:"recipient,omitempty"`
	Message   string     `json:"message,omitempty"`
	SendAt    *time.Time `json:"send_at,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newNotificationResponse(n notification.Notification) notificationResponse {
	resp := notificationResponse{
		ID:        n.ID,
		Status:    string(n.Status),
		Recipient: n.Recipient,
		Message:   n.Message,
	}
	if !n.SendAt.IsZero() {
		resp.SendAt = &n.SendAt
	}
	return resp
}

func writeError(c *ginext.Context, status int, message string) {
	c.JSON(status, errorResponse{Error: message})
}

func writeServiceError(c *ginext.Context, err error) {
	switch {
	case errors.Is(err, notification.ErrInvalidID),
		errors.Is(err, notification.ErrInvalidRecipient),
		errors.Is(err, notification.ErrInvalidMessage),
		errors.Is(err, notification.ErrInvalidSendAt):
		writeError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, notification.ErrNotFound):
		writeError(c, http.StatusNotFound, err.Error())
	case errors.Is(err, notification.ErrAlreadyCanceled),
		errors.Is(err, notification.ErrInvalidStatusTransition):
		writeError(c, http.StatusConflict, err.Error())
	case errors.Is(err, notification.ErrNotImplemented):
		writeError(c, http.StatusNotImplemented, err.Error())
	default:
		writeError(c, http.StatusInternalServerError, err.Error())
	}
}
