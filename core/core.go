package core

import (
	"sicsimgo/core/units"
)

var registers Registers = Registers{
	A:  units.Int24{},
	X:  units.Int24{},
	L:  units.Int24{},
	B:  units.Int24{},
	S:  units.Int24{},
	T:  units.Int24{},
	F:  units.Float48{},
	PC: units.Int24{},
	SW: units.Int24{},
}
var memory Memory = Memory{
	Data: make([]byte, MEMORY_SIZE),
}

type ExecuteState bool

const (
	ExecuteStartState ExecuteState = true
	ExecuteStopState  ExecuteState = false
)

var SimExecuteState ExecuteState = ExecuteStopState

// Public functions
func Reset() {
	SimExecuteState = ExecuteStopState
	registers = Registers{
		A:  units.Int24{},
		X:  units.Int24{},
		L:  units.Int24{},
		B:  units.Int24{},
		S:  units.Int24{},
		T:  units.Int24{},
		F:  units.Float48{},
		PC: units.Int24{},
		SW: units.Int24{},
	}
	memory = Memory{
		Data: make([]byte, MEMORY_SIZE),
	}
}
