package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/outgoing"
	"github.com/4aykovski/yadro_test_task/internal/model"
)

type TableService struct {
	openTime    time.Time
	closeTime   time.Time
	oneHourCost int

	arrivedClients map[string]struct{}
	clientQueue    []string

	tables     []model.Table
	busyTables map[string]*model.Table
}

func NewTableService(tables []model.Table, openTime, closeTime time.Time, oneHourCost int) *TableService {
	return &TableService{
		openTime:       openTime,
		closeTime:      closeTime,
		oneHourCost:    oneHourCost,
		arrivedClients: make(map[string]struct{}, len(tables)),
		clientQueue:    make([]string, 0, len(tables)),
		busyTables:     make(map[string]*model.Table, len(tables)),
		tables:         tables,
	}
}

type ClientArrivedDto struct {
	Time       time.Time
	ClientName string
}

func (s *TableService) ClientArrived(input ClientArrivedDto) event.Event {
	// проверить рабочие часы
	if !s.isOpenNow(input.Time) {
		errorEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrNotOpenYet)
		return errorEvent
	}

	// проверить присутствие клиента в пришедших
	if s.isClientInClub(input.ClientName) {
		errorEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrYouShallNotPass)
		return errorEvent
	}

	// добавить в пришедших клиентов
	s.arrivedClients[input.ClientName] = struct{}{}
	return nil
}

type ClientLeftDto struct {
	Time       time.Time
	ClientName string
}

func (s *TableService) ClientLeft(input ClientLeftDto) event.Event {
	// проверить в клубе ли клиент
	if !s.isClientInClub(input.ClientName) {
		outEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrClientUnknown)
		return outEvent
	}

	// посчитать прибыль и занимаемое время
	table, ok := s.busyTables[input.ClientName]
	if ok {
		table.Income += s.calculateIncome(*table, input.Time)
		table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, input.Time))
	}

	var freeTableId = 0
	// освободить стол, если клиент занимает стол
	if table, ok := s.busyTables[input.ClientName]; ok {
		table.IsTaken = false
		freeTableId = table.Id
		delete(s.busyTables, input.ClientName)
	}

	// удалить из клуба
	delete(s.arrivedClients, input.ClientName)

	// посадить клиента из очереди за освободившийся стол
	if len(s.clientQueue) > 0 && freeTableId != 0 {
		client := s.clientQueue[0]
		s.clientQueue = s.clientQueue[1:]
		s.changeClientTable(client, freeTableId, input.Time)

		outEvent := outgoing.NewClientTookPlaceEvent(input.Time, client, freeTableId)
		return outEvent
	}

	return nil
}

type ClientTookPlaceDto struct {
	Time       time.Time
	TableId    int
	ClientName string
}

func (s *TableService) ClientTookPlace(input ClientTookPlaceDto) event.Event {
	// проверить присутствие клиента в клубе
	if !s.isClientInClub(input.ClientName) {
		errorEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrClientUnknown)
		return errorEvent
	}

	// проверить занят ли стол за который садится клиент
	if !s.isTableFree(input.TableId) {
		errorEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrPlaceIsBusy)
		return errorEvent
	}

	// посчитать выручку, если пересаживается, и занимаемое время
	table, ok := s.busyTables[input.ClientName]
	if ok {
		table.Income += s.calculateIncome(*table, input.Time)
		table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, input.Time))
	}

	// поменять стол клиента
	s.changeClientTable(input.ClientName, input.TableId, input.Time)
	return nil
}

type ClientWaitingDto struct {
	Time       time.Time
	ClientName string
}

func (s *TableService) ClientWaiting(input ClientWaitingDto) event.Event {
	// проверить присутствие клиента в клубе
	if !s.isClientInClub(input.ClientName) {
		errorEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrClientUnknown)
		return errorEvent
	}

	// проверить свободные столы
	if s.isThereFreeTables() {
		errEvent := outgoing.NewErrorEvent(input.Time, outgoing.ErrICanWaitNoLonger)
		return errEvent
	}

	// проверить длину очереди
	if len(s.clientQueue) > len(s.tables) {
		errEvent := outgoing.NewClientLeftEvent(input.Time, input.ClientName)
		return errEvent
	}

	// поместить в очередь
	s.clientQueue = append(s.clientQueue, input.ClientName)
	return nil
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
	clients := make([]string, 0, len(s.arrivedClients))
	for clientName := range s.arrivedClients {
		clients = append(clients, clientName)
	}
	sort.Strings(clients)

	for _, clientName := range clients {
		if table, ok := s.busyTables[clientName]; ok {
			table.Income += s.calculateIncome(*table, s.closeTime)
			table.WasTakenFor = table.WasTakenFor.Add(s.calculateBusyTime(*table, s.closeTime))
			delete(s.busyTables, clientName)
		}

		delete(s.arrivedClients, clientName)

		outgoingEvent := outgoing.NewClientLeftEvent(s.closeTime, clientName)
		fmt.Println(outgoingEvent.String())
	}

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
	_, ok := s.arrivedClients[clientName]
	return ok
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

	var sub = roundedLeaveAt.Sub(roundedTakenAt)
	var diff time.Duration
	if int(sub.Minutes())%60 == 0 {
		diff = sub.Truncate(time.Hour)
	} else {
		diff = time.Hour + sub.Truncate(time.Hour)
	}

	return s.oneHourCost * int(diff.Hours())
}

func (s *TableService) calculateBusyTime(table model.Table, busyTime time.Time) time.Duration {
	return busyTime.Sub(table.TakenAt)
}
