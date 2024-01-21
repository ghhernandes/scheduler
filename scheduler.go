package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MatchType int

const (
	MatchPrefix MatchType = iota
	MatchRef
)

type ScheduleRef uuid.UUID

func (sr ScheduleRef) String() string {
	return fmt.Sprintf("schedule-%s", uuid.UUID(sr).String())
}

type Value struct {
	Date    time.Time
	Webhook string
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
