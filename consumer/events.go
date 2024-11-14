package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ntentasd/ecommerce-microservices/models"
)

type EventType string

var (
	ProductCreated = EventType("ProductCreated")
	OrderCreated   = EventType("OrderCreated")

	ErrEventTypeNotFound = errors.New("event type not found")
)

type Events map[EventType][]models.Event

type EventStore struct {
	data Events
	mu   sync.RWMutex
}

func (s *EventStore) Add(eventType EventType, event models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[eventType] = append(s.data[eventType], event)
}

func (s *EventStore) Get(eventType EventType, index int) (models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.data[eventType]
	if !ok {
		return nil, ErrEventTypeNotFound
	}

	event := events[index]

	return event, nil
}

func (s *EventStore) Pop(eventType EventType) (models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data[eventType]) == 0 {
		return nil, fmt.Errorf("no events available for type: %v", eventType)
	}

	event := s.data[eventType][0]
	s.data[eventType] = s.data[eventType][1:]

	return event, nil
}

func (s *EventStore) IsEmpty(eventType EventType) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data[eventType]) == 0
}
