package incoming

import (
	"fmt"
	"time"
)

type Type int

const (
	ClientArrived Type = iota + 1
	ClientTookPlace
	ClientWaiting
	ClientLeft
)

type Event struct {
	time       time.Time
	_type      Type
	clientName string
	tableId    int
}

func New(time time.Time, _type Type, clientName string, placeId int) *Event {
	switch _type {
	case ClientArrived:
		return NewClientArrived(time, clientName)
	case ClientTookPlace:
		return NewClientTookPlace(time, clientName, placeId)
	case ClientWaiting:
		return NewClientWaiting(time, clientName)
	case ClientLeft:
		return NewClientLeft(time, clientName)
	default:
		return nil
	}
}

func NewClientArrived(time time.Time, clientName string) *Event {
	return &Event{
		time:       time,
		_type:      ClientArrived,
		clientName: clientName,
	}
}

func NewClientTookPlace(time time.Time, clientName string, placeId int) *Event {
	return &Event{
		time:       time,
		_type:      ClientTookPlace,
		clientName: clientName,
		tableId:    placeId,
	}
}

func NewClientWaiting(time time.Time, clientName string) *Event {
	return &Event{
		time:       time,
		_type:      ClientWaiting,
		clientName: clientName,
	}
}

func NewClientLeft(time time.Time, clientName string) *Event {
	return &Event{
		time:       time,
		_type:      ClientLeft,
		clientName: clientName,
	}
}

func (e *Event) Type() int {
	return int(e._type)
}

func (e *Event) String() string {
	switch e._type {
	case ClientArrived, ClientLeft, ClientWaiting:
		return fmt.Sprintf("%v %v %v", e.time.Format("15:04"), e._type, e.clientName)
	case ClientTookPlace:
		return fmt.Sprintf("%v %v %v %v", e.time.Format("15:04"), e._type, e.clientName, e.tableId)
	default:
		return fmt.Sprintf("%v %v", e.time.Format("15:04"), e._type)
	}
}

func (e *Event) Time() time.Time {
	return e.time
}

func (e *Event) ClientName() string {
	return e.clientName
}

func (e *Event) TableId() int {
	return e.tableId
}
