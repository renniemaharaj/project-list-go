package conveyor

import (
	"time"

	"github.com/renniemaharaj/conveyor/pkg/conveyor"
)

var (
	singleton = conveyor.CreateManager().SetDebugging(false).
		SetMaxWorkers(100).SetMinWorkers(0).SetSafeQueueLength(10).
		SetTimePerTicker(10 / time.Second).Start()
)

// Get returns the shared conveyor belt
func Get() *conveyor.ConveyorBelt {
	return singleton.B
}
