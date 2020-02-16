package calendar

import (
	"testing"
	"time"

	"github.com/dark705/otus/hw11/internal/calendar/event"
	"github.com/dark705/otus/hw11/internal/storage"
	"github.com/sirupsen/logrus"
)

func TestNewCalendarHaveNoEvents(t *testing.T) {
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	events, err := calendar.GetAllEvents()
	if err != ErrNoEventsInStorage || len(events) != 0 {
		t.Error("In new storage exist events")
	}
}

func TestAddEventSuccess(t *testing.T) {
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	event, _ := event.CreateEvent("2006-01-02T15:00:00Z", "2006-01-02T16:00:00Z", "Event 1", "Some Desc1")
	err := calendar.AddEvent(event)
	if err != nil {
		t.Error("Can't add event to storage")
	}

	events, err := calendar.GetAllEvents()
	if err != nil || len(events) != 1 {
		t.Error("In storage not 1 event")
	}
}

func TestDelEventSuccess(t *testing.T) {
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	event, _ := event.CreateEvent("2006-01-02T15:00:00Z", "2006-01-02T16:00:00Z", "Event 1", "Some Desc1")
	_ = calendar.AddEvent(event)
	err := calendar.DelEvent(0)
	if err != nil {
		t.Error("Can't del event from storage")
	}

	events, err := calendar.GetAllEvents()
	if err == nil || len(events) != 0 {
		t.Error("In storage exist events")
	}
}

func TestAddDateIntervalBusy(t *testing.T) {
	var err error
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	event1, _ := event.CreateEvent("2006-01-02T15:00:00Z", "2006-01-02T16:00:00Z", "Event 1", "Some Desc1")
	event2, _ := event.CreateEvent("2006-01-02T16:00:00Z", "2006-01-02T17:00:00Z", "Event 2", "Some Desc2")
	event3, _ := event.CreateEvent("2006-01-02T18:00:00Z", "2006-01-02T19:00:00Z", "Event 3", "Some Desc3")
	err = calendar.AddEvent(event1)
	err = calendar.AddEvent(event2)
	err = calendar.AddEvent(event3)
	if err != nil {
		t.Error("Error on add not intersection events")
	}

	event4, _ := event.CreateEvent("2006-01-02T16:10:00Z", "2006-01-02T16:20:00Z", "Event 4", "Some Desc4")
	err = calendar.AddEvent(event4)
	if err != ErrDateBusy {
		t.Error("Add not return error for busy interval")
	}

	event5, _ := event.CreateEvent("2006-01-02T10:10:00Z", "2006-01-02T22:00:00Z", "Event 5", "Some Desc5")
	err = calendar.AddEvent(event5)
	if err != ErrDateBusy {
		t.Error("Add not return error for busy interval")
	}

	event6, _ := event.CreateEvent("2006-01-02T17:10:00Z", "2006-01-02T18:10:00Z", "Event 6", "Some Desc6")
	err = calendar.AddEvent(event6)
	if err != ErrDateBusy {
		t.Error("Add not return error for busy interval")
	}
}

func TestGetEvent(t *testing.T) {
	var err error
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	event, _ := event.CreateEvent("2006-01-02T15:00:00Z", "2006-01-02T16:00:00Z", "Event 1", "Some Desc1")
	_ = calendar.AddEvent(event)

	_, err = calendar.GetEvent(1)

	if err == nil {
		t.Error("Not get error, for not exist event")
	}
	getEvent, err := calendar.GetEvent(0)
	if err != nil {
		t.Error("Get error, for exist event")
	}

	if getEvent.StartTime != event.StartTime ||
		getEvent.EndTime != event.EndTime ||
		getEvent.Title != event.Title ||
		getEvent.Description != event.Description {
		t.Error("Event in storage not ident")
	}
}

func TestEditEvent(t *testing.T) {
	inMemory := storage.InMemory{}
	inMemory.Init()
	calendar := Calendar{Storage: &inMemory, Logger: &logrus.Logger{}}

	event, _ := event.CreateEvent("2006-01-02T15:00:00Z", "2006-01-02T16:00:00Z", "Event 1", "Some Desc1")
	_ = calendar.AddEvent(event)

	editEvent, _ := calendar.GetEvent(0)
	editEvent.StartTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:10:00Z")
	editEvent.EndTime, _ = time.Parse(time.RFC3339, "2006-01-02T15:20:00Z")
	editEvent.Title = "newTitle"
	editEvent.Description = "newDescription"

	err := calendar.EditEvent(editEvent)
	if err != nil {
		t.Error("Got not expected error on edit")
	}

	eventFromStorageAfterEdit, _ := calendar.GetEvent(0)
	if eventFromStorageAfterEdit != editEvent {
		t.Error("Edit Event not ident Event in storage after edit")
	}
}
