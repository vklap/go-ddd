package event_handlers

import (
	"context"
	"fmt"
	"github.com/vklap/go_ddd/internal/adapters"
	"github.com/vklap/go_ddd/internal/domain/command_model"
	"github.com/vklap/go_ddd/pkg/ddd"
)

// EmailChangedEventHandler implements ddd.EventHandler.
type EmailChangedEventHandler struct {
	pubSubClient adapters.PubSubClient
	events       []ddd.Event
}

// NewEmailChangedEventHandler is a constructor function to be used by the Bootstrapper.
func NewEmailChangedEventHandler(pubSubClient adapters.PubSubClient) *EmailChangedEventHandler {
	return &EmailChangedEventHandler{pubSubClient: pubSubClient, events: make([]ddd.Event, 0)}
}

// Handle manages the business logic flow, and is the glue between the Domain and the Adapters.
func (h *EmailChangedEventHandler) Handle(ctx context.Context, event ddd.Event) error {
	e, ok := event.(*command_model.EmailChangedEvent)
	if ok == false {
		panic(fmt.Sprintf("failed to handle email changed: want %T, got %T", &command_model.NotifySlackEvent{}, e))
	}
	if err := h.pubSubClient.NotifyEmailChanged(ctx, e.UserID, e.NewEmail, e.OriginalEmail); err != nil {
		return err
	}
	h.events = append(h.events, &command_model.NotifySlackEvent{Message: "Email Notification sent for EmailChanged event"})
	return nil
}

// Commit is responsible for committing the changes performed by the Handle method, such as
// committing a database transaction managed by the repository.
// This method is being called by the framework, so it should not be called from within the Handle method.
func (h *EmailChangedEventHandler) Commit(ctx context.Context) error {
	return h.pubSubClient.Commit(ctx)
}

// Rollback is responsible to rollback changes performed by the Handle method, such as
// rollback a database transaction managed by the repository.
// This method is being called by the framework, so it should not be called from within the Handle method.
func (h *EmailChangedEventHandler) Rollback(ctx context.Context) error {
	return h.pubSubClient.Rollback(ctx)
}

// Events reports about events.
// These events will be handled by the DDD framework if appropriate event handlers were registered by the bootstrapper.
// This method is being called by the framework, so it should not be called from within the Handle method.
func (h *EmailChangedEventHandler) Events() []ddd.Event {
	return h.events
}

var _ ddd.EventHandler = (*EmailChangedEventHandler)(nil)
