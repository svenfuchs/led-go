# Led

A line editor in Go. Inspired by [linenoise](https://github.com/antirez/linenoise),
but written with extensibility and separation of concerns in mind.

The motivation behind this library is to provide line editor functionality that
can be controlled, extended, and modified by clients with regards to how to
handle key events. In this regard the name of the library also is a pun on
letting go (of control).

## Usage

```go
import (
	"github.com/svenfuchs/led"
)

func main() {
	NewLed("$ ").Run()
}
```

See [example/led.go](/blob/master/example/led.go) for a usage example that makes
use of custom key handlers, suggestions, completion, and history, and reimplements
(most of?) the functionality in linenoise.

Also see https://github.com/svenfuchs/travis-go for an example that takes over
more control.

## Todo

* Handle multiple lines if line length exceeds terminal width

