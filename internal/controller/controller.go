package controller

import (
	"fmt"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
)

type tableService interface {
	ClientArrived(event event.Event) string
	ClientLeft(event event.Event) string
	ClientTookPlace(event event.Event) string
	ClientWaiting(event event.Event) string
	OpenTime() time.Time
	CloseTime() time.Time
	PrintIncome()
	Close()
}

type Controller struct {
	tableService tableService
}

func New(tableService tableService) *Controller {
	return &Controller{
		tableService: tableService,
	}
}

func (c *Controller) HandleEvent(event event.Event) (string, error) {
	switch event.Type() {
	case incoming.ClientArrived:
		return c.tableService.ClientArrived(event), nil
	case incoming.ClientLeft:
		return c.tableService.ClientLeft(event), nil
	case incoming.ClientTookPlace:
		return c.tableService.ClientTookPlace(event), nil
	case incoming.ClientWaiting:
		return c.tableService.ClientWaiting(event), nil
	}
	return "", fmt.Errorf("unknown event type: %v", event.Type())
}

func (c *Controller) Close() {
	c.tableService.Close()
}

func (c *Controller) PrintIncome() {
	c.tableService.PrintIncome()
}

func (c *Controller) OpenTime() time.Time {
	return c.tableService.OpenTime()
}

func (c *Controller) CloseTime() time.Time {
	return c.tableService.CloseTime()
}
