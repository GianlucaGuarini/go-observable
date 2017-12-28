# go-observable

[![Build Status](https://travis-ci.org/GianlucaGuarini/go-observable.svg?branch=master)](https://travis-ci.org/GianlucaGuarini/go-observable)
[![Coverage Status](https://coveralls.io/repos/github/GianlucaGuarini/go-observable/badge.svg?branch=master)](https://coveralls.io/github/GianlucaGuarini/go-observable?branch=master)

Golang script heavily inspired to [riot-observable](https://github.com/riot/observable)

It allows you to send and receive events with a tiny simple API (thread safe)

## Installation

```go
go get github.com/GianlucaGuarini/go-observable
```

## Api

### New()

Create a new observable struct reference

```go
o := observable.New()
```

### On(event string, fn interface{})

Subscribe a callback to a certain event key

```go
o.On("ready", func() {
  // I am ready!
})
```

You can subscribe your observable also to multiple events for example:

```go
o.On("stop error", func(args ...interface{}) {
  // It will be called when the "error" or the "stop" event will be dispatched
  // args[0] here is either "stop" or "error" depending on the Trigger call
})
```

By using the `*` namespace you will be able to listen all the observable events

__IMPORTANT__<br/>
If you pass arguments to your listeners you should add them also to this listener
as well

```go
o.On("*", func(args ...interface{}) {
  // args[0] here will be the event key called in the Trigger method
  // in this case "stop" and "start"
})
o.Trigger("stop").Trigger("start")
```



### Off(event string, fn interface{})

Unsubscribe a callback from an event key

```go
onReady := func() {
  // I am ready
}
o.On("ready", onReady)
// do your stuff...
o.Off("ready", onReady) // the onReady will not be called anymore
```

By using the key `*` you will remove all the observed

```go
o.On("ready", func() {
  // ready
})
o.On("stop", func() {
  // stop
})

o.Off("*") // kill all the listeners

o.Trigger("ready stop") // the callbacks above will never be called
```

### One(event string, fn interface{})

Subscribe a callback in order to be called only once

```go
o.One("ready", func() {
  // I am ready and I will not be called anymore
})
```

### Trigger(event string, arguments ...interface{})

Call all the callbacks subscribed to a certain event

```go
o.On("message", func(message string) {
  // do something with the message
})
o.Trigger("message", "Hello There!")
```



