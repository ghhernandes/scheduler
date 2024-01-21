package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type MatchType int

const (
	MatchPrefix MatchType = iota
	MatchRef
)

type ScheduleRef uuid.UUID

func NewRef() ScheduleRef {
	return ScheduleRef(uuid.New())
}

func NewRefString(key string) ScheduleRef {
	return ScheduleRef(uuid.Must(uuid.Parse(key)))
}

func (sr ScheduleRef) String() string {
	return uuid.UUID(sr).String()
}

type Value struct {
	Date    time.Time `json:"date"`
	Webhook string    `json:"webhook"`
}

type ValueVer struct {
	Version int
	Value   Value
}

type Schedule struct {
	Ref   ScheduleRef
	Value Value

	Prev *ValueVer
}

type Matcher struct {
	Type  MatchType
	Value string
}
