package service

import (
	"testing"
	"time"
)

func TestCalendarService_CreateEvent(t *testing.T) {
	service := NewCalendarService()

	tests := []struct {
		name     string
		userID   int
		date     time.Time
		title    string
		wantErr  bool
	}{
		{
			name:    "valid event",
			userID:  1,
			date:    time.Now(),
			title:   "Meeting",
			wantErr: false,
		},
		{
			name:    "empty title",
			userID:  1,
			date:    time.Now(),
			title:   "",
			wantErr: true,
		},
		{
			name:    "zero date",
			userID:  1,
			date:    time.Time{},
			title:   "Event",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateEvent(tt.userID, tt.date, tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarService_GetDayEvents(t *testing.T) {
	service := NewCalendarService()
	now := time.Now()
	userID := 1

	service.CreateEvent(userID, now, "Event 1")
	service.CreateEvent(userID, now.Add(24*time.Hour), "Event 2")
	service.CreateEvent(2, now, "Event 3")

	events := service.GetDayEvents(userID, now)
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Title != "Event 1" {
		t.Errorf("Expected 'Event 1', got '%s'", events[0].Title)
	}
}

func TestCalendarService_UpdateEvent(t *testing.T) {
	service := NewCalendarService()
	userID := 1
	date := time.Now()

	event, _ := service.CreateEvent(userID, date, "Original Title")

	updatedEvent, err := service.UpdateEvent(event.ID, userID, date, "Updated Title")
	if err != nil {
		t.Errorf("UpdateEvent() unexpected error: %v", err)
	}

	if updatedEvent.Title != "Updated Title" {
		t.Errorf("Expected 'Updated Title', got '%s'", updatedEvent.Title)
	}

	_, err = service.UpdateEvent("non-existent", userID, date, "Title")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	service := NewCalendarService()
	userID := 1
	date := time.Now()

	event, _ := service.CreateEvent(userID, date, "Event to delete")

	err := service.DeleteEvent(event.ID)
	if err != nil {
		t.Errorf("DeleteEvent() unexpected error: %v", err)
	}

	events := service.GetDayEvents(userID, date)
	if len(events) != 0 {
		t.Errorf("Expected 0 events after deletion, got %d", len(events))
	}

	err = service.DeleteEvent("non-existent")
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}