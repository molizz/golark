package golark

import (
	"github.com/molizz/golark/callback"
	"github.com/molizz/golark/utils"
)

var eventHub *Event

func init() {
	eventHub = &Event{
		hub: make(map[string]callback.EventProcessor),
	}
}

func RegisterEventProcessor(procs ...callback.EventProcessor) {
	if len(procs) == 0 {
		return
	}
	for _, p := range procs {
		eventHub.Register(p)
	}
}

type Event struct {
	hub map[string]callback.EventProcessor
}

func (e *Event) Register(proc callback.EventProcessor) {
	e.hub[proc.TypeLabel()] = proc
}

func (e *Event) Find(typeName string) callback.EventProcessor {
	proc, ok := e.hub[typeName]
	if !ok {
		utils.DefaultLog.Printf("lark callback '%s', not found processor\n", typeName)
		return nil
	}
	return proc
}
