package observable

import (
  "reflect"
  "strings"
  "sync"
)

// Helpers

// Add a callback under a certain event namespace
func (o *Observable) addCallback(event string, cb interface{}, isUnique bool) *Observable {
  fn := reflect.ValueOf(cb)

  events := strings.Fields(event)

  for _, s := range events {
    o.Lock()
    // does this namespace already exist?
    if !o.hasEvent(s) {
      o.Callbacks[s] = make([]callback, 1)
      o.Callbacks[s][0] = callback{fn, isUnique, false}
    } else {
      o.Callbacks[s] = append(o.Callbacks[s], callback{fn, isUnique, false})
    }
    o.Unlock()
  }

  return o
}

// check whether the Observable struct has already registered the event namespace
func (o *Observable) hasEvent(event string) bool {
  _, ok := o.Callbacks[event]
  return ok
}

// Structs

// private struct
type callback struct {
  fn        reflect.Value
  isUnique  bool
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
  events := strings.Fields(event)

  for key, param := range params {
    arguments[key] = reflect.ValueOf(param)
  }

  for _, s := range events {
    // check if the observable has already created this events map
    if o.hasEvent(s) {

      // loop all the callbacks
      // avoiding to call twice the ones registered with Observable.One
      for i, cb := range o.Callbacks[s] {
        if !cb.isUnique || cb.isUnique && !cb.wasCalled {
          cb.fn.Call(arguments)
        }
        // kill the callbacks registered with one
        if cb.isUnique {
          o.Off(s, o.Callbacks[s][i])
        }
        o.Callbacks[s][i].wasCalled = true
      }
    }
  }

  return o
}

// Off - stop listening a particular event
func (o *Observable) Off(event string, fn interface{}) *Observable {

  // try to get the value of the function we want unsubscribe
  fn = reflect.ValueOf(fn)

  // loop all the callbacks registered under the event namespace
  for i, cb := range o.Callbacks[event] {
    if fn == cb.fn {
      o.Lock()
      o.Callbacks[event] = append(o.Callbacks[event][:i], o.Callbacks[event][i+1:]...)
      o.Unlock()
    }
  }

  // if there are no more callbacks using this namespace
  // delete the key from the map
  if len(o.Callbacks[event]) == 0 {
    delete(o.Callbacks, event)
  }

  return o
}

// One - call the callback only once
func (o *Observable) One(event string, cb interface{}) *Observable {
  return o.addCallback(event, cb, true)
}
