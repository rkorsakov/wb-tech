package service

import (
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	ErrNotFound   = errors.New("event not found")
	ErrInvalidDate = errors.New("invalid date")
)

type Event struct {
	ID     string    `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

type CalendarService struct {
	mu     sync.RWMutex
	events map[string]*Event
}

func NewCalendarService() *CalendarService {
	return &CalendarService{
		events: make(map[string]*Event),
	}
}

func (s *CalendarService) CreateEvent(userID int, date time.Time, title string) (*Event, error) {
	if title == "" {
		return nil, errors.New("event title can't be blank")
	}
	if date.IsZero() {
		return nil, ErrInvalidDate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event := &Event{
		ID:     uuid.New().String(),
		UserID: userID,
		Date:   date,
		Title:  title,
	}
	s.events[event.ID] = event
	return event, nil
}

func (s *CalendarService) UpdateEvent(eventID string, userID int, date time.Time, title string) (*Event, error) {
	if eventID == "" {
		return nil, errors.New("event ID is required")
	}
	if title == "" {
		return nil, errors.New("event title can't be blank")
	}
	if date.IsZero() {
		return nil, ErrInvalidDate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existingEvent, exists := s.events[eventID]
	if !exists {
		return nil, ErrNotFound
	}

	if existingEvent.UserID != userID {
		return nil, errors.New("cannot change event owner")
	}

	existingEvent.Date = date
	existingEvent.Title = title

	return existingEvent, nil
}

func (s *CalendarService) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[id]; exists {
		delete(s.events, id)
	} else {
		return ErrNotFound
	}
	return nil
}

func (s *CalendarService) GetDayEvents(userID int, date time.Time) []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]*Event, 0)
	for _, event := range s.events {
		if event.UserID == userID && isSameDay(event.Date, date) {
			events = append(events, event)
		}
	}
	return events
}

func (s *CalendarService) GetWeekEvents(userID int, date time.Time) []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]*Event, 0)
	targetYear, targetWeek := date.ISOWeek()
	for _, event := range s.events {
		eventYear, eventWeek := event.Date.ISOWeek()
		if event.UserID == userID && eventYear == targetYear && eventWeek == targetWeek {
			events = append(events, event)
		}
	}
	return events
}

func (s *CalendarService) GetMonthEvents(userID int, date time.Time) []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]*Event, 0)
	targetYear, targetMonth := date.Year(), date.Month()
	for _, event := range s.events {
		if event.UserID == userID {
			eventYear, eventMonth := event.Date.Year(), event.Date.Month()
			if eventYear == targetYear && eventMonth == targetMonth {
				events = append(events, event)
			}
		}
	}
	return events
}

func isSameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}