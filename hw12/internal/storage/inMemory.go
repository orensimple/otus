package storage

import (
	"errors"
	"fmt"

	"github.com/dark705/otus/hw12/internal/calendar/event"
)

type InMemory struct {
	Events map[int]event.Event
	LastId int
}

func (s *InMemory) Init() error {
	s.Events = make(map[int]event.Event)
	return nil
}

func (s *InMemory) Add(e event.Event) error {
	e.Id = len(s.Events)
	s.Events[len(s.Events)] = e
	return nil
}

func (s *InMemory) Del(id int) error {
	delete(s.Events, id)
	return nil
}

func (s *InMemory) Get(id int) (event.Event, error) {
	event, exist := s.Events[id]
	if !exist {
		return event, errors.New(fmt.Sprintf("Event with id: %d not found", id))
	}
	return event, nil
}

func (s *InMemory) GetAll() ([]event.Event, error) {
	events := make([]event.Event, 0, len(s.Events))
	for _, e := range s.Events {
		events = append(events, e)
	}
	return events, nil
}

func (s *InMemory) Edit(e event.Event) error {
	_, exist := s.Events[e.Id]
	if !exist {
		return errors.New(fmt.Sprintf("Event with id: %d not found", e.Id))
	}
	s.Events[e.Id] = e
	return nil
}

func (s *InMemory) IntervalIsBusy(newEvent event.Event, new bool) (bool, error) {
	for id, existEvent := range s.Events {
		if newEvent.Id == id && new == false {
			continue
		}
		if existEvent.StartTime.Before(newEvent.EndTime) && existEvent.EndTime.After(newEvent.StartTime) {
			return true, nil
		}
	}
	return false, nil
}
