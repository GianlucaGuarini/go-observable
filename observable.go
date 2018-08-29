package observable

import (
	"reflect"
	"sync"
)

// event key uset to listen and remove all the events
const ALL_EVENTS_NAMESPACE = "*"

// Structs

// private struct
type callback struct {
	fn        reflect.Value
	isUnique  bool
	isTyped   bool
	wasCalled bool
}

// Public Observable struct
type Observable struct {
	Callbacks map[string][]callback
	*sync.RWMutex
}

// Public API

// New - returns a new observable reference
func New() *Observable {
	return &Observable{
		make(map[string][]callback),
		&sync.RWMutex{},
	}
}

// On - adds a callback function
func (o *Observable) On(event string, cb interface{}) *Observable {
	return o.addCallback(event, cb, false)
}

// Trigger - a particular event passing custom arguments
func (o *Observable) Trigger(event string, params ...interface{}) *Observable {
	o.Lock()
	defer o.Unlock()

	// get the arguments we want to pass to our listeners callbaks
	arguments := make([]reflect.Value, len(params))

	// get all the arguments
	for i, param := range params {
		arguments[i] = reflect.ValueOf(param)
	}

	o.dispatchEvent(event, arguments)

	// trigger the all events callback whenever this event was defined
	if o.hasEvent(ALL_EVENTS_NAMESPACE) && event != ALL_EVENTS_NAMESPACE {
		o.dispatchEvent(ALL_EVENTS_NAMESPACE, append([]reflect.Value{reflect.ValueOf(event)}, arguments...))
	}

	return o
}

// Off - stop listening a particular event
func (o *Observable) Off(event string, args ...interface{}) *Observable {
	o.Lock()
	defer o.Unlock()
	return o.offNoSync(event, args...)
}

/// Ditto, necessary for internal calls
func (o *Observable) offNoSync(event string, args ...interface{}) *Observable {
	if len(args) == 0 {
		// wipe all the event listeners
		if event == ALL_EVENTS_NAMESPACE {
			o.Callbacks = make(map[string][]callback)
		}
	} else if len(args) == 1 {
		o.removeEvent(event, args[0])
	} else {
		panic("Multiple off callbacks are not supported")
	}

	return o
}

// One - call the callback only once
func (o *Observable) One(event string, cb interface{}) *Observable {
	return o.addCallback(event, cb, true)
}
