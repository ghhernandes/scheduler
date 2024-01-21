package web

import (
	"time"

	"github.com/ghhernandes/scheduler"
)

type ScheduleVerResult struct {
	Version int    `json:"version"`
	Value   string `json:"value"`
}

type ValueResult struct {
	Date    time.Time `json:"date"`
	Webhook string    `json:"webhook"`
}

type ScheduleResult struct {
	ValueResult
	Ref      string              `json:"ref"`
	Versions []ScheduleVerResult `json:"versions"`
}

type ScheduleCreateRequest struct {
	Date    time.Time `json:"date"`
	Webhook string    `json:"webhook"`
}

func (scr ScheduleCreateRequest) ToValue() scheduler.Value {
	return scheduler.Value{
		Date:    scr.Date,
		Webhook: scr.Webhook,
	}
}

type ScheduleCreateResult struct {
	Ref string `json:"ref"`
}
