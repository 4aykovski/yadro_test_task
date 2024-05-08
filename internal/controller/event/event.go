package event

import (
	"github.com/4aykovski/yadro_test_task/internal/controller/event/incoming"
)

type Event interface {
	Type() incoming.Type
	String() string
}
