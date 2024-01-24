package lock

import "github.com/ghhernandes/scheduler"

type Locker interface {
	Lock(ref scheduler.ScheduleRef) error
	Unlock(ref scheduler.ScheduleRef) error
}
