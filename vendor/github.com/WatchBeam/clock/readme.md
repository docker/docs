# clock [![GoDoc](https://godoc.org/github.com/WatchBeam/clock?status.svg)](https://godoc.org/github.com/WatchBeam/clock) [![Build Status](https://travis-ci.org/WatchBeam/clock.svg)](https://travis-ci.org/WatchBeam/clock)

Time utility with lovely mocking support.

This is essentially a replacement for the `time` package which allows you to seamlessly swap in mock times, timers, and tickers. See the godocs (link above) for more detailed usage.

### Example

**hello.go**

```go
package main

import (
    "fmt"
    "github.com/WatchBeam/clock"
)

func main() {
    fmt.Printf("the time is %s", displayer{clock.C}.formatted())
}

type displayer struct {
    c clock.Clock
}

func (d displayer) formatted() string {
    now := d.c.Now()
    return fmt.Sprintf("%d:%d:%d", now.Hour(), now.Minute(), now.Second())
}
```

**hello_test.go**

```go
package main

import (
    "testing"
    "time"

    "github.com/WatchBeam/clock"
    "github.com/stretchr/testify/assert"
)

func TestDisplaysCorrectly(t *testing.T) {
    date, _ := time.Parse(time.UnixDate, "Sat Mar  7 11:12:39 PST 2015")
    c := clock.NewMockClock(date)
    d := displayer{c}

    assert.Equal(t, "11:12:39", d.formatted())
    c.AddTime(42 * time.Second)
    assert.Equal(t, "11:13:21", d.formatted())
}
```

### API & Compatibility

The API provided by this package and the mock version is nearly identical to that of the `time` package, with two notable differences:
 - The channel for Ticker and Timer instances an accessed via the `.Chan()` method, rather than reading the `.C` property. This allows the structures to be swapped out for their mock variants.
 - The mock Ticker never skips ticks when time advances. This allows you to call `.AddTime`/`.SetTime` on the mock clock without having to advance to each "ticked" time.
