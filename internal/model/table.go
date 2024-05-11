package model

import "time"

type Table struct {
	Id          int
	IsTaken     bool
	TakenAt     time.Time
	WasTakenFor time.Time
	Income      int
}
