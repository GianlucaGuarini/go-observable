# go-observable

Golang library heavily inspired to [riot-observable](https://github.com/riot/observable)

It allows to send and receive events with a tiny simple API

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

### One(event string, fn interface{})

Subscribe a callback in order to be called only once

```go
o.One("ready", func(){
  // I am ready and I will not be called anymore
})
```

### Trigger(event string, arguments ...interface{})

Call all the callbacks subscribed to a certain event

```go

o.On("message", func(message string){
  // do something with the message
})

o.Trigger("message", "Hello There!")

```



