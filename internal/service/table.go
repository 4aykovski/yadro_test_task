package service

import (
	"fmt"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/outgoing"
	"github.com/4aykovski/yadro_test_task/internal/model"
)

type TableService struct {
	openTime    time.Time
	closeTime   time.Time
	oneHourCost int

	arrivedClients []string
	clientQueue    []string

	tables     []model.Table
	busyTables map[string]*model.Table
}

func NewTableService(tables []model.Table, openTime, closeTime time.Time, oneHourCost int) *TableService {
	return &TableService{
		openTime:       openTime,
		closeTime:      closeTime,
		oneHourCost:    oneHourCost,
		arrivedClients: make([]string, 0, len(tables)),
		clientQueue:    make([]string, 0, len(tables)),
		busyTables:     make(map[string]*model.Table, len(tables)),
		tables:         tables,
	}
}

func (s *TableService) ClientArrived(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить рабочие часы
	if !s.isOpenNow(inEvent.Time()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrNotOpenYet)
		return errorEvent.String()
	}

	// проверить присутствие клиента в пришедших
	if s.isClientInClub(inEvent.ClientName()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrYouShallNotPass)
		return errorEvent.String()
	}

	// добавить в пришедших клиентов
	s.arrivedClients = append(s.arrivedClients, inEvent.ClientName())
	return ""
}

func (s *TableService) ClientLeft(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить в клубе ли клиент
	if !s.isClientInClub(inEvent.ClientName()) {
		outEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrClientUnknown)
		return outEvent.String()
	}

	// посчитать прибыль и занимаемое время
	table, ok := s.busyTables[inEvent.ClientName()]
	if ok {
		table.Income += s.calculateIncome(*table, inEvent.Time())
		table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, inEvent.Time()))
	}

	// освободить стол
	s.busyTables[inEvent.ClientName()].IsTaken = false
	freeTableId := s.busyTables[inEvent.ClientName()].Id
	delete(s.busyTables, inEvent.ClientName())

	// удалить из клуба
	for i, client := range s.arrivedClients {
		if client == inEvent.ClientName() {
			s.arrivedClients = append(s.arrivedClients[:i], s.arrivedClients[i+1:]...)
			break
		}
	}

	// посадить клиента из очереди за освободившийся стол
	if len(s.clientQueue) > 0 {
		client := s.clientQueue[0]
		s.clientQueue = s.clientQueue[1:]
		s.changeClientTable(client, freeTableId, inEvent.Time())

		outEvent := outgoing.NewClientTookPlaceEvent(inEvent.Time(), client, freeTableId)
		return outEvent.String()
	}

	return ""
}

func (s *TableService) ClientTookPlace(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить присутствие клиента в клубе
	if !s.isClientInClub(inEvent.ClientName()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrClientUnknown)
		return errorEvent.String()
	}

	// проверить занят ли стол за который садится клиент
	if !s.isTableFree(inEvent.TableId()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrPlaceIsBusy)
		return errorEvent.String()
	}

	// посчитать выручку, если пересаживается, и занимаемое время
	table, ok := s.busyTables[inEvent.ClientName()]
	if ok {
		table.Income += s.calculateIncome(*table, inEvent.Time())
		table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, inEvent.Time()))
	}

	// поменять стол клиента
	s.changeClientTable(inEvent.ClientName(), inEvent.TableId(), inEvent.Time())
	return ""
}

func (s *TableService) ClientWaiting(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить свободные столы
	if s.isThereFreeTables() {
		errEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrICanWaitNoLonger)
		return errEvent.String()
	}

	// проверить длину очереди
	if len(s.clientQueue) > len(s.tables) {
		errEvent := outgoing.NewClientLeftEvent(inEvent.Time(), inEvent.ClientName())
		return errEvent.String()
	}

	// поместить в очередь
	s.clientQueue = append(s.clientQueue, inEvent.ClientName())
	return ""
}

func (s *TableService) OpenTime() time.Time {
	return s.openTime
}

func (s *TableService) CloseTime() time.Time {
	return s.closeTime
}

func (s *TableService) PrintIncome() {
	for _, table := range s.tables {
		fmt.Printf("%v %v %v\n", table.Id, table.Income, table.WasTakenFor.Format("15:04"))
	}
}

func (s *TableService) Close() {
	// выгнать всех
	for clientName, table := range s.busyTables {
		table.Income += s.calculateIncome(*table, s.closeTime)
		table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, s.closeTime))

		outgoingEvent := outgoing.NewClientLeftEvent(s.closeTime, clientName)
		table.TakenAt = time.Time{}
		table.IsTaken = false

		fmt.Println(outgoingEvent.String())
	}
}

func (s *TableService) isOpenNow(time time.Time) bool {
	return time.After(s.openTime) && time.Before(s.closeTime)
}

func (s *TableService) isClientInClub(clientName string) bool {
	for _, client := range s.arrivedClients {
		if client == clientName {
			return true
		}
	}
	return false
}

func (s *TableService) isTableFree(tableId int) bool {
	for _, table := range s.tables {
		if table.Id == tableId && !table.IsTaken {
			return true
		}
	}
	return false
}

func (s *TableService) changeClientTable(clientId string, tableId int, takenTime time.Time) {
	table, ok := s.busyTables[clientId]
	if ok {
		table.IsTaken = false
		table.TakenAt = time.Time{}
		delete(s.busyTables, clientId)
	}

	table = &s.tables[tableId-1]
	table.IsTaken = true
	table.TakenAt = takenTime

	s.busyTables[clientId] = table
}

func (s *TableService) isThereFreeTables() bool {
	for _, table := range s.tables {
		if !table.IsTaken {
			return true
		}
	}
	return false
}

func (s *TableService) calculateIncome(table model.Table, leaveTime time.Time) int {
	roundedTakenAt := table.TakenAt
	roundedLeaveAt := leaveTime
	diff := time.Hour + roundedLeaveAt.Sub(roundedTakenAt).Truncate(time.Hour)
	return s.oneHourCost * int(diff.Hours())
}

func (s *TableService) calculateBusyTime(table model.Table, busyTime time.Time) time.Duration {
	return busyTime.Sub(table.TakenAt)
}
