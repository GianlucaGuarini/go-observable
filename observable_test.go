package observable_test

import (
  "github.com/GianlucaGuarini/go-observable"
  "sync"
  "testing"
  "time"
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

func TestOnTriggerMultipleEvensString(t *testing.T) {
  o := observable.New()
  n := 0

  var lastEvt string

  o.On("foo bar", func(eventName string) {
    lastEvt = eventName
    n++
  })

  o.Trigger("foo bar").Trigger("bar foo")

  if lastEvt != "foo" {
    t.Errorf("The last event name triggered is %s instead of being %s", lastEvt, "foo")
  }

  if n != 4 {
    t.Errorf("The counter is %d instead of being %d", n, 4)
  }
}

func TestOffMultipleEvensString(t *testing.T) {
  o := observable.New()
  n := 0

	increment := func(args ...interface{}) {
		n++
	}

  o.On("foo bar", increment).On("baz", increment)

  o.Off("foo bar baz", increment)

  o.Trigger("foo bar baz")

  if n != 0 {
    t.Errorf("The counter is %d instead of being %d", n, 0)
  }

}

func TestOnAll(t *testing.T) {
  o := observable.New()
  n := 0
  var lastEvt string

  onAll := func(eventName string) {
    lastEvt = eventName
    n++
  }

  o.On("*", onAll)

  o.Trigger("foo bar").Trigger("foo").Trigger("bar")

  o.Off("*", onAll)

  o.Trigger("foo bar").Trigger("foo")

  if lastEvt != "bar" {
    t.Errorf("The last event name triggered is %s instead of being %s", lastEvt, "bar")
  }

  if n != 3 {
    t.Errorf("The counter is %d instead of being %d", n, 3)
  }

}

func TestOffAll(t *testing.T) {
  o := observable.New()
  n := 0

  o.On("foo", func() {
    n++
  })

  o.On("bar", func() {
    n++
  })

  o.Off("*")

  o.Trigger("foo").Trigger("bar").Trigger("foo bar")

  if n != 0 {
    t.Errorf("The counter is %d instead of being %d", n, 0)
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

func TestRace(t *testing.T) {

  o := observable.New()
  n := 0

  asyncTask := func(wg *sync.WaitGroup) {
    o.Trigger("foo")
    wg.Done()
  }
  var wg sync.WaitGroup

  wg.Add(5)

  o.On("foo", func() {
    n++
  })

  go asyncTask(&wg)
  go asyncTask(&wg)
  go asyncTask(&wg)
  go asyncTask(&wg)
  go asyncTask(&wg)

  wg.Wait()

  if n != 5 {
    t.Errorf("The counter is %d instead of being %d", n, 5)
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

/**
Test using `On` / `defer Off` and `Trigger` concurrently

One useful pattern (for SSE, Server Send Event) is to have code code that does `Trigger`
events and a server handler that turns the event `On` / `Off` on disconnection.
However this can lead to a race condition as `Trigger` and `On`/`Off` are not currently
synchronized.
*/
func TestOnOffTriggerConcurrency(b *testing.T) {
	o := observable.New()

	waiter := make(chan struct{})

	// This is called synchronously
	observerFunc := func(args ...interface{}) {
		time.Sleep(7 * time.Millisecond)
	}

	// Our "http listener" accepts many connections
	// This test can accidentally pass, but it cannot accidentally fail
	go func() {
		<-waiter
		for i := 0; i < 100; i++ {
			o.On("baguette", observerFunc)
			time.Sleep(10 * time.Millisecond)
			o.Off("baguette")
		}
	}()

	// Does not exists yet
	o.Trigger("baguette")
	waiter <- struct{}{}
	for i := 0; i < 150; i++ {
		o.Trigger("baguette")
	}
}
