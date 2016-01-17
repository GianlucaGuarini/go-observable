package observable

import (
  "reflect"
)

type callback struct {
  fn        reflect.Value
  isOne     bool
  wasCalled bool
}

// Observable struct
type Observable struct {
  Callbacks   map[string][]callback
  argumentsCh map[string]chan []reflect.Value
  doneCh      map[string]chan int
}

// New - returns a observable struct
func New() *Observable {
  return &Observable{
    make(map[string][]callback),
    make(map[string]chan []reflect.Value),
    make(map[string]chan int),
  }
}

// On - adds a callback function
func (o *Observable) On(event string, fn interface{}) *Observable {
  o.addCallback(event, fn, false)
  return o
}

// Trigger - a particular event passing custom arguments
func (o *Observable) Trigger(event string, params ...interface{}) *Observable {
  arguments := make([]reflect.Value, len(params))
  for key, param := range params {
    arguments[key] = reflect.ValueOf(param)
  }
  o.argumentsCh[event] <- arguments
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
    o.doneCh[event] <- 1
    delete(o.Callbacks, event)
    delete(o.doneCh, event)
    delete(o.argumentsCh, event)
  }

  return o
}

// One - call the callback only once
func (o *Observable) One(event string, fn interface{}) *Observable {
  o.addCallback(event, fn, true)
  return o
}
