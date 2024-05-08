package outgoing

import (
	"fmt"
	"time"
)

type Type int

const (
	ClientLeft = iota + 11
	ClientTookPlace
	Error
)

const (
	ErrNotOpenYet       = "NotOpenYet"
	ErrICanWaitNoLonger = "ICanWaitNoLonger"
	ErrYouShallNotPass  = "YouShallNotPass"
	ErrPlaceIsBusy      = "PlaceIsBusy"
	ErrClientUnknown    = "ClientUnknown"
)

type Event struct {
	time       time.Time
	_type      Type
	clientName string
	msg        string
	tableId    int
}

func NewClientLeftEvent(time time.Time, clientName string) *Event {
	return &Event{
		time:       time,
		_type:      ClientLeft,
		clientName: clientName,
	}
}

func NewClientTookPlaceEvent(time time.Time, clientName string, placeId int) *Event {
	return &Event{
		time:       time,
		_type:      ClientTookPlace,
		clientName: clientName,
		tableId:    placeId,
	}
}

func NewErrorEvent(time time.Time, msg string) *Event {
	return &Event{
		time:  time,
		_type: Error,
		msg:   msg,
	}
}

func (e *Event) Type() Type {
	return e._type
}

func (e *Event) String() string {
	switch e.Type() {
	case ClientLeft:
		return fmt.Sprintf("%v %v %v", e.time.Format("15:04"), e._type, e.clientName)
	case ClientTookPlace:
		return fmt.Sprintf("%v %v %v %v", e.time.Format("15:04"), e._type, e.clientName, e.tableId)
	case Error:
		return fmt.Sprintf("%v %v %v", e.time.Format("15:04"), e._type, e.msg)
	default:
		return fmt.Sprintf("%v %v", e.time.Format("15:04"), e._type)
	}
}
