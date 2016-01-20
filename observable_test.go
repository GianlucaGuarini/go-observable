package observable_test

import (
  "github.com/GianlucaGuarini/go-observable"
  "testing"
)

func TestOn(t *testing.T) {

  o := observable.New()
  n := 0

  o.On("foo", func() {
    n++
  }).On("bar", func() {
    n++
  }).On("foo", func() {
    n++
  })

  o.Trigger("foo").Trigger("foo").Trigger("bar")

  if n != 5 {
    t.Errorf("The counter is %d instead of being %d", n, 5)
  }

}

func TestOff(t *testing.T) {
  o := observable.New()
  n := 0

  onFoo1 := func() {
    n++
  }

  onFoo2 := func() {
    n++
  }

  o.On("foo", onFoo1).On("foo", onFoo2)

  o.Off("foo", onFoo1).Off("foo", onFoo2).On("foo", onFoo1)

  o.Trigger("foo")

  if n != 1 {
    t.Errorf("The counter is %d instead of being %d", n, 1)
  }

}

func TestOne(t *testing.T) {
  o := observable.New()
  n := 0

  onFoo := func() {
    n++
  }

  o.One("foo", onFoo)

  o.Trigger("foo").Trigger("foo").Trigger("foo")

  if n != 1 {
    t.Errorf("The counter is %d instead of being %d", n, 1)
  }

}

func TestArguments(t *testing.T) {
  o := observable.New()
  n := 0
  o.On("foo", func(arg1 bool, arg2 string) {
    n++
    if arg1 != true || arg2 != "bar" {
      t.Error("The arguments must be correctly passed to the callback")
    }
  })

  o.Trigger("foo", true, "bar")

  if n != 1 {
    t.Errorf("The counter is %d instead of being %d", n, 1)
  }
}

func TestTrigger(t *testing.T) {
  o := observable.New()
  // the trigger without any listener should not throw errors
  o.Trigger("foo")
}

/**
 * Speed Benchmarks
 */

var eventsList = []string{"foo", "bar", "baz", "boo"}

func BenchmarkOnTrigger(b *testing.B) {
  o := observable.New()
  n := 0

  for _, e := range eventsList {
    o.On(e, func() {
      n++
    })
  }

  for i := 0; i < b.N; i++ {
    for _, e := range eventsList {
      o.Trigger(e)
    }
  }
}
