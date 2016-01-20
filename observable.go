package observable

import "reflect"

// Helpers

func (o *Observable) addCallback(event string, fn interface{}, isOne bool) {
  if !o.hasEvent(event) {
    o.Callbacks[event] = make([]callback, 1)
    o.Callbacks[event][0] = callback{reflect.ValueOf(fn), isOne, false}
  } else {
    o.Callbacks[event] = append(o.Callbacks[event], callback{reflect.ValueOf(fn), isOne, false})
  }
}

func (o *Observable) hasEvent(event string) bool {
  _, ok := o.Callbacks[event]
  return ok
}

// Structs

type callback struct {
  fn        reflect.Value
  isOne     bool
  wasCalled bool
}

// Observable struct
type Observable struct {
  Callbacks map[string][]callback
}

// Public API

// New - returns a observable struct
func New() *Observable {
  return &Observable{
    make(map[string][]callback),
  }
}

// On - adds a callback function
func (o *Observable) On(event string, fn interface{}) *Observable {
  o.addCallback(event, fn, false)
  return o
}

// Trigger - a particular event passing custom arguments
func (o *Observable) Trigger(event string, params ...interface{}) *Observable {

  // check if the observable has already created this events map
  if o.hasEvent(event) {
    arguments := make([]reflect.Value, len(params))
    for key, param := range params {
      arguments[key] = reflect.ValueOf(param)
    }

    for i, cb := range o.Callbacks[event] {
      if cb.isOne && !cb.wasCalled || !cb.isOne {
        cb.fn.Call(arguments)
      }
      if cb.isOne {
        o.Off(event, o.Callbacks[event][i])
      }
      o.Callbacks[event][i].wasCalled = true
    }
  }

  return o
}

// Off - stop listening a particular event
func (o *Observable) Off(event string, fn interface{}) *Observable {

  fn = reflect.ValueOf(fn)

  for i, cb := range o.Callbacks[event] {
    if fn == cb.fn {
      o.Callbacks[event] = append(o.Callbacks[event][:i], o.Callbacks[event][i+1:]...)
    }
  }

  if len(o.Callbacks[event]) == 0 {
    delete(o.Callbacks, event)
  }

  return o
}

// One - call the callback only once
func (o *Observable) One(event string, fn interface{}) *Observable {
  o.addCallback(event, fn, true)
  return o
}
