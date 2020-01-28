package Storage

import (
	"github.com/dark705/otus/hw08/internal/Calendar/Event"
)

type InMemory struct {
	Events []Event.Event
}

func (s *InMemory) Add(e Event.Event) error {
	if s.intervalIsBusy(e) {
		return NewErrDateBusy()
	}
	s.Events = append(s.Events, e)
	return nil
}

func (s *InMemory) Del(id int) error {
	if id >= len(s.Events) {
		return NewErrNotFoundWithId(id)
	}
	//fast, but change order
	s.Events[id] = s.Events[len(s.Events)-1]
	s.Events = s.Events[:len(s.Events)-1]
	return nil
}

func (s *InMemory) intervalIsBusy(newEvent Event.Event) bool {
	for _, existEvent := range s.Events {
		//NewEvent include existEvent
		if newEvent.StartTime.Before(existEvent.StartTime) && newEvent.EndTime.After(existEvent.EndTime) {
			return true
		}
		//EndTime of newEvent inside existEvent
		if newEvent.EndTime.After(existEvent.StartTime) && newEvent.EndTime.Before(existEvent.EndTime) {
			return true
		}
		//StartTime of newEvent inside existEvent
		if newEvent.StartTime.After(existEvent.StartTime) && newEvent.StartTime.Before(existEvent.EndTime) {
			return true
		}

	}
	return false
}
