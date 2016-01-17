package observable

import (
  "reflect"
)

func (o *Observable) listen(argumentsCh chan []reflect.Value, doneCh chan int, event string) {
  for {
    select {
    case params := <-argumentsCh:
      for i, cb := range o.Callbacks[event] {
        if cb.isOne && !cb.wasCalled || !cb.isOne {
          cb.fn.Call(params)
        }
        if cb.isOne {
          o.Off(event, o.Callbacks[event][i])
        }
        o.Callbacks[event][i].wasCalled = true
      }
    case <-doneCh:
      // kill the loop
      return
    }

  }
}

func (o *Observable) addCallback(event string, fn interface{}, isOne bool) {

  if !o.hasEvent(event) {

    o.Callbacks[event] = make([]callback, 1)
    o.argumentsCh[event] = make(chan []reflect.Value)
    o.doneCh[event] = make(chan int)
    o.Callbacks[event][0] = callback{reflect.ValueOf(fn), isOne, false}

    go o.listen(o.argumentsCh[event], o.doneCh[event], event)

  } else {
    o.Callbacks[event] = append(o.Callbacks[event], callback{reflect.ValueOf(fn), isOne, false})
  }
}

func (o *Observable) hasEvent(event string) bool {
  _, ok := o.Callbacks[event]
  return ok
}
