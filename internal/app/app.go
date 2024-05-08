package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/4aykovski/yadro_test_task/internal/controller"
	"github.com/4aykovski/yadro_test_task/internal/controller/event"
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
	"github.com/4aykovski/yadro_test_task/internal/model"
	"github.com/4aykovski/yadro_test_task/pkg/helpers"
)

type Controller interface {
	HandleEvent(event event.Event) (string, error)
	OpenTime() time.Time
	CloseTime() time.Time
	PrintIncome()
	Close()
}

type System struct {
	events     []event.Event
	controller Controller
}

func New(events []event.Event, controller Controller) *System {
	return &System{
		events:     events,
		controller: controller,
	}
}

func (s *System) Run() error {
	fmt.Println(s.controller.OpenTime().Format("15:04"))

	for _, e := range s.events {
		fmt.Println(e.String())
		res, err := s.controller.HandleEvent(e)
		if err != nil {
			return err
		}

		if res != "" {
			fmt.Println(res)
		}
	}

	s.controller.Close()

	fmt.Println(s.controller.CloseTime().Format("15:04"))
	s.controller.PrintIncome()
	return nil
}

func Run(data []string) error {
	tablesCount, err := strconv.Atoi(data[0])
	if err != nil {
		return fmt.Errorf("can't parse tables count: %w", err)
	}

	var tables []model.Table
	for i := 0; i < tablesCount; i++ {
		tables = append(tables, model.Table{
			Id:      i + 1,
			IsTaken: false,
		})
	}

	splitTime := strings.Split(data[1], " ")
	openTime, err := helpers.ParseTime(splitTime[0])
	if err != nil {
		return fmt.Errorf("can't parse time: %w", err)
	}
	closeTime, err := helpers.ParseTime(splitTime[1])
	if err != nil {
		return fmt.Errorf("can't parse time: %w", err)
	}

	oneHourCost, err := strconv.Atoi(data[2])
	if err != nil {
		return fmt.Errorf("can't parse tables count: %w", err)
	}

	events, err := ParseEvents(data[3:])
	if err != nil {
		return fmt.Errorf("can't parse events: %w", err)
	}

	cont := controller.New(tables, openTime, closeTime, oneHourCost)

	sys := New(events, cont)

	return sys.Run()
}

func ParseEvents(data []string) ([]event.Event, error) {
	events := make([]event.Event, 0, len(data))

	for _, line := range data {
		split := strings.Split(line, " ")

		eventTime, err := helpers.ParseTime(split[0])
		if err != nil {
			return nil, fmt.Errorf("can't parse time")
		}

		_type, ok := helpers.ParsePositiveInt(split[1])
		if !ok {
			return nil, fmt.Errorf("can't parse event type")
		}

		name := split[2]

		table := 0
		if len(split) == 4 {
			table, ok = helpers.ParsePositiveInt(split[3])
			if !ok {
				return nil, fmt.Errorf("can't parse table")
			}
		}

		events = append(events, incoming.New(eventTime, incoming.Type(_type), name, table))
	}
	return events, nil
}
