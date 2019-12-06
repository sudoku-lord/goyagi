# logger-go

Lob's logger for Go.

## Usage

Include the logger-go package in your import.

```go
import (
    "github.com/lob/logger-go"
)
```

Logs will always include:

* logging level
* hostname
* release
  * Defined by the environment variable `RELEASE`
  * This is set by Convox in staging/production environments
* timestamp
* nanoseconds

Logs be written to stdout by default. The package ships with a global logger so you can easily emit logs without having to instantiate a logger. This default logger will always write to stdout.

```go
logger.Info("Hello, world!", logger.Data{"fun": "things"})
// Outputs: {"level":"info","host":"Kyle.local","release":"test12345","data":{"fun":"things"},"nanoseconds":1532024420744842400,"timestamp":"2018-07-19T11:20:20-07:00","message":"Hello, world!"}
```

Alternatively, you can instantiate your own logger. This is useful if you want to attach additional data or top-level information to your logger, which will force all logs emitted by that logger to include that info.

```go
l1 := logger.New().ID("test")
l2 := l1.Data(logger.Data{"test": "data"})

l1.Info("hi")
// Outputs {"level":"info","host":"HOSTNAME","release":"RELEASE","id":"test","nanoseconds":1531945897647586415,"timestamp":"2018-07-18T13:31:37-07:00","message":"hi"}
l2.Info("hi")
// Outputs {"level":"info","host":"HOSTNAME","release":"RELEASE","id":"test","data":{"test":"data"},"nanoseconds":1531945897647593709,"timestamp":"2018-07-18T13:31:37-07:00","message":"hi"}

// If Data or Root are empty, they will not show up in the logs.
l1 = l1.Data(logger.Data{})
l1.Info("hi")
// Outputs {"level":"info","host":"HOSTNAME","release":"RELEASE","id":"test","nanoseconds":1531945897647586415,"timestamp":"2018-07-18T13:31:37-07:00","message":"hi"}

// To log errors, use Err. If it's a normal error, a runtime stack trace is logged. This provides limited context, so it's recommended to use pkg/errors instead (see below).
err := fmt.Errorf("foo")
l1.Err(err).Error("unknown error")
// {"level":"error","host":"HOSTNAME","release":"RELEASE","id":"test","error":{"message":"foo","stack":"goroutine 1 [running]:\ngithub.com/lob/logger-go.Logger.log(0x111b0c0, 0xc420010440, 0x0, 0x0, 0x0, 0xc4200b8200, 0x19, 0x1f4, 0xc420010450, 0x1, ...)\n\t/go/src/github.com/lob/logger-go/logger.go:153 +0x5d2\ngithub.com/lob/logger-go.Logger.Error(0x111b0c0, 0xc420010440, 0x0, 0x0, 0x0, 0xc4200b8200, 0x19, 0x1f4, 0xc420010450, 0x1, ...)\n\t/go/src/github.com/lob/logger-go/logger.go:101 +0xce\nmain.main()\n\t/go/src/github.com/lob/logger-go/main.go:27 +0x5db\n"},"nanoseconds":1531945897647586415,"timestamp":"2018-07-18T13:31:37-07:00","message":"unknown error"}

// If the error is wrapped with pkg/errors, a better stack trace is logged. See https://godoc.org/github.com/pkg/errors#hdr-Retrieving_the_stack_trace_of_an_error_or_wrapper for more info.
err = errors.New("bar")
l1.Err(err).Error("unknown error")
// {"level":"error","host":"HOSTNAME","release":"RELEASE","id":"test","error":{"message":"bar","stack":"\nmain.main\n\t/go/src/github.com/lob/logger-go/main.go:26\nruntime.main\n\t/.goenv/versions/1.10.3/src/runtime/proc.go:198\nruntime.goexit\n\t/.goenv/versions/1.10.3/src/runtime/asm_amd64.s:2361"},"nanoseconds":1531945897647586415,"timestamp":"2018-07-18T13:31:37-07:00","message":"unknown error"}
```

If you want the logger to use a specific writer to pipe logs to anywhere other than stdout, please use the `NewWithWriter` method. Make sure the argument you are passing in implements the writer interface.
```go
type CustomWriter struct{}
func (cw *CustomWriter) Write(b []byte) (int, error) {
	// your custom write logic here
}

loggerWithWriter := logger.NewWithWriter(CustomWriter{})
```

The logger supports five levels of logging.

### Info

```go
logger.Info("Hello, world!")
```

### Error

```go
logger.Error("Hello, world!")
```

### Warn

```go
logger.Warn("Hello, world!")
```

### Debug

```go
logger.Debug("Hello, world!")
```

### Fatal

```go
logger.Fatal("Hello, world!")
```

We currently do not support trace-level logging, since zerolog, the underlying logging library, does not support trace (or custom level logging).

## Echo Middleware

This package also comes with middleware to be used with the
[Echo](https://github.com/labstack/echo) web framework. To use it, you can just
register it with `e.Use()`:

```go
e := echo.New()

e.Use(logger.Middleware())
```

There are also some scenarios where you don't wany an error to be logged and
registered with Echo. An example is for broken pipe errors. When this happens,
the client closed the connection, so even though it manifests as a network
error, it's not actionable for us, so we would rather ignore it. To do that, you
can use `logger.MiddlewareWithConfig()` and `logger.MiddlewareConfig`.

```go
e := echo.New()

e.Use(logger.MiddlewareWithConfig(logger.MiddlewareConfig{
	IsIgnorableError: func(err error) bool {
		e := errors.Cause(err)

		if netErr, ok := e.(*net.OpError); ok {
			if osErr, ok := netErr.Err.(*os.SyscallError); ok {
				return osErr.Err.Error() == syscall.EPIPE.Error() || osErr.Err.Error() == syscall.ECONNRESET.Error()
			}
		}

		return false
	},
}))
```

With this middleware, not only will it create a `handled request` log line for
every request, but it will also attach a request-specific logger to the Echo
context. To use it, you can pull it out with `logger.FromEchoContext()`:

```go
e.GET("/", func(ctx echo.Context) error {
        log := logger.FromEchoContext(ctx)

        log.Info("in handler")

        // ...

        return nil
})
```

**You should always try to use this logger instead of creating your own so that
it contains contextual information about the request.**

## Errors and Stack Traces

In Go, the idiomatic `error` type doesn't contain any stack information by default. Since there's no mechanism to extract a stack, it's common to use [`runtime.Stack`](https://golang.org/pkg/runtime/#Stack) to generate one. The problem with that is since the stack is usually created at the point of _error handling_ and not the point of _error creation_, the stack most likely won't have the origination point of the error.

Because of this, the community has created a way to maintain stack information as the error gets bubbled up while still adhering to the idiomatic `error` interface. This is with the [`github.com/pkg/errors` package](https://godoc.org/github.com/pkg/errors). It allows developers to add context and stack frames to errors that are generated throughout a codebase. By leveraging this, you can produce a much better stack trace.

To demonstate the difference between these two mechanisms, here is an example:

```go
func main() {
	log := logger.New()

	nativeErr := nativeFunction1()
	pkgErr := pkgFunction1()

	log.Err(nativeErr).Error("native error")
	log.Err(pkgErr).Error("pkg error")
}

func nativeFunction1() error {
	return nativeFunction2() // line 21
}

func nativeFunction2() error {
	// returns a native error
	return fmt.Errorf("foo") // line 26
}

func pkgFunction1() error {
	return pkgFunction2() // line 30
}

func pkgFunction2() error {
	// returns a pkg/errors error
	return errors.New("foo") // line 35
}
```

This sample code produces these stack traces:

```
goroutine 1 [running]:
github.com/lob/logger-go.Logger.log(0x111a8c0, 0xc420010440, 0x0, 0x0, 0x0, 0xc4200b8200, 0x19, 0x1f4, 0xc420010450, 0x1, ...)
        /Users/robinjoseph/go/src/github.com/lob/logger-go/logger.go:154 +0x5ad
github.com/lob/logger-go.Logger.Error(0x111a8c0, 0xc420010440, 0x0, 0x0, 0x0, 0xc4200b8200, 0x19, 0x1f4, 0xc420010450, 0x1, ...)
        /Users/robinjoseph/go/src/github.com/lob/logger-go/logger.go:101 +0xce
main.main()
        /Users/robinjoseph/go/src/github.com/lob/logger-go/ex/main.go:16 +0x1bd


main.pkgFunction2
        /Users/robinjoseph/go/src/github.com/lob/logger-go/ex/main.go:35
main.pkgFunction1
        /Users/robinjoseph/go/src/github.com/lob/logger-go/ex/main.go:30
main.main
        /Users/robinjoseph/go/src/github.com/lob/logger-go/ex/main.go:14
runtime.main
        /Users/robinjoseph/.goenv/versions/1.10.3/src/runtime/proc.go:198
runtime.goexit
        /Users/robinjoseph/.goenv/versions/1.10.3/src/runtime/asm_amd64.s:2361
```

As you can see from the runtime stack (the former), it contains neither function names nor line numbers to indicate the original cause of the error. But with the `pkg/errors` stack trace (the latter), it's very clear and contains all the necessary information needed for debugging.

This logging package will attempt to extract the `pkg/errors` stack trace if that information exists, but otherwise, it will provide the runtime stack. **It's because of this that we strongly recommend adding `pkg/errors` to wrap all errors in your codebase.** This will aid in the general debuggability of your applications.

For more information on how to wrap errors, like [`errors.WithStack()`](https://godoc.org/github.com/pkg/errors#WithStack), check out the [`pkg/errors` docs](https://godoc.org/github.com/pkg/errors#hdr-Adding_context_to_an_error).

## Development

```sh
# Install necessary dependencies for development
make setup

# Ensure dependencies are safe, complete, and reproducible
make deps

# Run tests and generate test coverage report
make test

# Run linter
make lint

# Remove temporary files generated by the build command and the build directory
make clean
```

## Cutting a New Release

After new commits have been added, a new git tag should also be created so that
tools like `go mod` and `dep` can utilize it for more semantic versioning.

If there is a breaking change introduced in this new version, then it should be
a major version bump. If there are no breaking changes and only new features,
then it should be a minor version bump. And if there are no breaking changes, no
new features, and only bug fixes, then it should be a patch version.

After determining what the next version should be, make sure you're on an
up-to-date master branch and run `make release`. **The `v` in the version tag is
important for tools to recognize it.**

```sh
$ git checkout master
$ git fetch
$ git rebase
$ make release tag=v1.0.0 # replace 1.0.0 with the next version
```

## Help

The following command generates a list of all `make` commands.

```sh
make help
```

## FAQ

### Tweaking code coverage

See this [blog post about coverage in Go](https://blog.golang.org/cover).

### Tweaking the linter

See the [gometalinter documentation](https://github.com/alecthomas/gometalinter#installing).
