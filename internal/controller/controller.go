package controller

import (
	"fmt"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
	"github.com/4aykovski/yadro_test_task/internal/service"
)

type tableService interface {
	ClientArrived(input service.ClientArrivedDto) string
	ClientLeft(input service.ClientLeftDto) string
	ClientTookPlace(input service.ClientTookPlaceDto) string
	ClientWaiting(input service.ClientWaitingDto) string
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
	inEvent, ok := event.(*incoming.Event)
	if !ok {
		return "", fmt.Errorf("unexpected event: %v", event)
	}

	switch event.Type() {
	case incoming.ClientArrived:
		dto := service.ClientArrivedDto{
			ClientName: inEvent.ClientName(),
			Time:       inEvent.Time(),
		}

		return c.tableService.ClientArrived(dto), nil
	case incoming.ClientLeft:
		dto := service.ClientLeftDto{
			ClientName: inEvent.ClientName(),
			Time:       inEvent.Time(),
		}
		return c.tableService.ClientLeft(dto), nil
	case incoming.ClientTookPlace:
		dto := service.ClientTookPlaceDto{
			ClientName: inEvent.ClientName(),
			TableId:    inEvent.TableId(),
			Time:       inEvent.Time(),
		}
		return c.tableService.ClientTookPlace(dto), nil
	case incoming.ClientWaiting:
		dto := service.ClientWaitingDto{
			ClientName: inEvent.ClientName(),
			Time:       inEvent.Time(),
		}

		return c.tableService.ClientWaiting(dto), nil
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
