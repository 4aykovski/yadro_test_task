package controller

import (
	"fmt"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/outgoing"
	"github.com/4aykovski/yadro_test_task/internal/model"
)

type Controller struct {
	openTime    time.Time
	closeTime   time.Time
	oneHourCost int

	arrivedClients []string
	clientQueue    []string

	tables     []model.Table
	busyTables map[string]*model.Table
}

func New(tables []model.Table, openTime, closeTime time.Time, oneHourCost int) *Controller {
	return &Controller{
		openTime:       openTime,
		closeTime:      closeTime,
		oneHourCost:    oneHourCost,
		arrivedClients: make([]string, 0, len(tables)),
		clientQueue:    make([]string, 0, len(tables)),
		busyTables:     make(map[string]*model.Table, len(tables)),
		tables:         tables,
	}
}

func (c *Controller) HandleEvent(event event.Event) (string, error) {
	switch event.Type() {
	case incoming.ClientArrived:
		return c.handleClientArrived(event), nil
	case incoming.ClientLeft:
		return c.handleClientLeft(event), nil
	case incoming.ClientTookPlace:
		return c.handleClientTookPlace(event), nil
	case incoming.ClientWaiting:
		return c.handleClientWaiting(event), nil
	}
	return "", fmt.Errorf("unknown event type: %v", event.Type())
}

func (c *Controller) OpenTime() time.Time {
	return c.openTime
}

func (c *Controller) CloseTime() time.Time {
	return c.closeTime
}

func (c *Controller) PrintIncome() {
	for _, table := range c.tables {
		fmt.Printf("%v %v %v\n", table.Id, table.Income, table.WasTakenFor.Format("15:04"))
	}
}

func (c *Controller) Close() {
	// выгнать всех
	for clientName, table := range c.busyTables {
		table.Income += c.calculateIncome(*table, c.closeTime)
		table.WasTakenFor = table.WasTakenFor.Add(c.calculateBusyTime(*table, c.closeTime))

		outgoingEvent := outgoing.NewClientLeftEvent(c.closeTime, clientName)
		table.TakenAt = time.Time{}
		table.IsTaken = false

		fmt.Println(outgoingEvent.String())
	}
}

func (c *Controller) handleClientArrived(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить рабочие часы
	if !c.isOpenNow(inEvent.Time()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrNotOpenYet)
		return errorEvent.String()
	}

	// проверить присутствие клиента в пришедших
	if c.isClientInClub(inEvent.ClientName()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrYouShallNotPass)
		return errorEvent.String()
	}

	// добавить в пришедших клиентов
	c.arrivedClients = append(c.arrivedClients, inEvent.ClientName())
	return ""
}

func (c *Controller) handleClientLeft(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить в клубе ли клиент
	if !c.isClientInClub(inEvent.ClientName()) {
		outEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrClientUnknown)
		return outEvent.String()
	}

	// посчитать прибыль и занимаемое время
	table, ok := c.busyTables[inEvent.ClientName()]
	if ok {
		table.Income += c.calculateIncome(*table, inEvent.Time())
		table.WasTakenFor = table.WasTakenFor.Add(c.calculateBusyTime(*table, inEvent.Time()))
	}

	// освободить стол
	c.busyTables[inEvent.ClientName()].IsTaken = false
	freeTableId := c.busyTables[inEvent.ClientName()].Id
	delete(c.busyTables, inEvent.ClientName())

	// удалить из клуба
	for i, client := range c.arrivedClients {
		if client == inEvent.ClientName() {
			c.arrivedClients = append(c.arrivedClients[:i], c.arrivedClients[i+1:]...)
			break
		}
	}

	// посадить клиента из очереди за освободившийся стол
	if len(c.clientQueue) > 0 {
		client := c.clientQueue[0]
		c.clientQueue = c.clientQueue[1:]
		c.changeClientTable(client, freeTableId, inEvent.Time())

		outEvent := outgoing.NewClientTookPlaceEvent(inEvent.Time(), client, freeTableId)
		return outEvent.String()
	}

	return ""
}

func (c *Controller) handleClientTookPlace(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить присутствие клиента в клубе
	if !c.isClientInClub(inEvent.ClientName()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrClientUnknown)
		return errorEvent.String()
	}

	// проверить занят ли стол за который садится клиент
	if !c.isTableFree(inEvent.TableId()) {
		errorEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrPlaceIsBusy)
		return errorEvent.String()
	}

	// посчитать выручку, если пересаживается, и занимаемое время
	table, ok := c.busyTables[inEvent.ClientName()]
	if ok {
		table.Income += c.calculateIncome(*table, inEvent.Time())
		table.WasTakenFor = table.WasTakenFor.Add(c.calculateBusyTime(*table, inEvent.Time()))
	}

	// поменять стол клиента
	c.changeClientTable(inEvent.ClientName(), inEvent.TableId(), inEvent.Time())
	return ""
}

func (c *Controller) handleClientWaiting(event event.Event) string {
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return fmt.Sprintf("unknown event: %v", event)
	}

	// проверить свободные столы
	if c.isThereFreeTables() {
		errEvent := outgoing.NewErrorEvent(inEvent.Time(), outgoing.ErrICanWaitNoLonger)
		return errEvent.String()
	}

	// проверить длину очереди
	if len(c.clientQueue) > len(c.tables) {
		errEvent := outgoing.NewClientLeftEvent(inEvent.Time(), inEvent.ClientName())
		return errEvent.String()
	}

	// поместить в очередь
	c.clientQueue = append(c.clientQueue, inEvent.ClientName())
	return ""
}

func (c *Controller) isOpenNow(time time.Time) bool {
	return time.After(c.openTime) && time.Before(c.closeTime)
}

func (c *Controller) isClientInClub(clientName string) bool {
	for _, client := range c.arrivedClients {
		if client == clientName {
			return true
		}
	}
	return false
}

func (c *Controller) isTableFree(tableId int) bool {
	for _, table := range c.tables {
		if table.Id == tableId && !table.IsTaken {
			return true
		}
	}
	return false
}

func (c *Controller) changeClientTable(clientId string, tableId int, takenTime time.Time) {
	table, ok := c.busyTables[clientId]
	if ok {
		table.IsTaken = false
		table.TakenAt = time.Time{}
		delete(c.busyTables, clientId)
	}

	table = &c.tables[tableId-1]
	table.IsTaken = true
	table.TakenAt = takenTime

	c.busyTables[clientId] = table
}

func (c *Controller) isThereFreeTables() bool {
	for _, table := range c.tables {
		if !table.IsTaken {
			return true
		}
	}
	return false
}

func (c *Controller) calculateIncome(table model.Table, leaveTime time.Time) int {
	roundedTakenAt := table.TakenAt
	roundedLeaveAt := leaveTime
	diff := time.Hour + roundedLeaveAt.Sub(roundedTakenAt).Truncate(time.Hour)
	return c.oneHourCost * int(diff.Hours())
}

func (c *Controller) calculateBusyTime(table model.Table, busyTime time.Time) time.Duration {
	return busyTime.Sub(table.TakenAt)
}
