# Stacktrace [![Circle CI](https://img.shields.io/circleci/project/palantir/stacktrace/master.svg?label=circleci)](https://circleci.com/gh/palantir/stacktrace) [![Travis CI](https://img.shields.io/travis/palantir/stacktrace/master.svg?label=travis)](https://travis-ci.org/palantir/stacktrace)

Look at Palantir, such a Java shop. I can't believe they want stack traces in
their Go code.

### Why would anyone want stack traces in Go code?

This is difficult to debug:

```
Inverse tachyon pulse failed
```

This gives the full story and is easier to debug:

```
Failed to register for villain discovery
 --- at github.com/palantir/shield/agent/discovery.go:265 (ShieldAgent.reallyRegister) ---
 --- at github.com/palantir/shield/connector/impl.go:89 (Connector.Register) ---
Caused by: Failed to load S.H.I.E.L.D. config from /opt/shield/conf/shield.yaml
 --- at github.com/palantir/shield/connector/config.go:44 (withShieldConfig) ---
Caused by: There isn't enough time (4 picoseconds required)
 --- at github.com/palantir/shield/axiom/pseudo/resource.go:46 (PseudoResource.Adjust) ---
 --- at github.com/palantir/shield/axiom/pseudo/growth.go:110 (reciprocatingPseudo.growDown) ---
 --- at github.com/palantir/shield/axiom/pseudo/growth.go:121 (reciprocatingPseudo.verify) ---
Caused by: Inverse tachyon pulse failed
 --- at github.com/palantir/shield/metaphysic/tachyon.go:72 (TryPulse) ---
```

Note that stack traces are *not designed to be user-visible*. We have found them
to be valuable in log files of server applications. Nobody wants to see these in
CLI output or a web interface or a return value from library code.

## Intent

The intent is *not* that we capture the exact state of the stack when an error
happens, including every function call. For a library that does that, see
[github.com/go-errors/errors](https://github.com/go-errors/errors). The intent
here is to attach relevant contextual information (messages, variables) at
strategic places along the call stack, keeping stack traces compact and
maximally useful.

## Example Usage

<!-- pre instead of code block to support bold text inside -->
<pre>
func WriteAll(baseDir string, entities []Entity) error {
    err := os.MkdirAll(baseDir, 0755)
    if err != nil {
        return <b>stacktrace.Propagate(err, "Failed to create base directory")</b>
    }
    for _, ent := range entities {
        path := filepath.Join(baseDir, fileNameForEntity(ent))
        err = Write(path, ent)
        if err != nil {
            return <b>stacktrace.Propagate(err, "Failed to write %v to %s", ent, path)</b>
        }
    }
    return nil
}
</pre>

## Functions

#### stacktrace.Propagate(cause error, msg string, vals ...interface{}) error

Propagate wraps an error to include line number information. This is going to be
your most common stacktrace call.

As in all of these functions, the `msg` and `vals` work like `fmt.Errorf`.

The message passed to Propagate should describe the action that failed,
resulting in `cause`. The canonical call looks like this:

<pre>
result, err := process(arg)
if err != nil {
    return nil, <b>stacktrace.Propagate(err, "Failed to process %v", arg)</b>
}
</pre>

To write the message, ask yourself "what does this call do?" What does
`process(arg)` do? It processes ${arg}, so the message is that we failed to
process ${arg}.

Pay attention that the message is not redundant with the one in `err`. In the
`WriteAll` example above, any error from `os.MkdirAll` will already contain the
path it failed to create, so it would be redundant to include it again in our
message. However, the error from `os.MkdirAll` will not identify that path as
corresponding to the "base directory" so we propagate with that information.

If it is not possible to add any useful contextual information beyond what is
already included in an error, `msg` can be an empty string:

<pre>
func Something() error {
    mutex.Lock()
    defer mutex.Unlock()

    err := reallySomething()
    return <b>stacktrace.Propagate(err, "")</b>
}
</pre>

The purpose of `""` as opposed to a separate function is to make you feel a
little guilty every time you do this.

This example also illustrates the behavior of Propagate when `cause` is nil
&ndash; it returns nil as well. There is no need to check `if err != nil`.

#### stacktrace.NewError(msg string, vals ...interface{}) error

NewError is a drop-in replacement for `fmt.Errorf` that includes line number
information. The canonical call looks like this:

<pre>
if !IsOkay(arg) {
    return <b>stacktrace.NewError("Expected %v to be okay", arg)</b>
}
</pre>

### Error Codes

Occasionally it can be useful to propagate an error code while unwinding the
stack. For example, a RESTful API may use the error code to set the HTTP status
code.

The type `stacktrace.ErrorCode` is a typedef for uint16. You name the set of
error codes relevant to your application.

```
const (
    EcodeManifestNotFound = stacktrace.ErrorCode(iota)
    EcodeBadInput
    EcodeTimeout
)
```

The special value `stacktrace.NoCode` is equal to `math.MaxUint16`, so avoid
using that. NoCode is the error code of errors with no code explicitly attached.

An ordinary `stacktrace.Propagate` preserves the error code of an error.

#### stacktrace.PropagateWithCode(cause error, code ErrorCode, msg string, vals ...interface{}) error

#### stacktrace.NewErrorWithCode(code ErrorCode, msg string, vals ...interface{}) error

PropagateWithCode and NewErrorWithCode are analogous to Propagate and NewError
but also attach an error code.

<pre>
_, err := os.Stat(manifestPath)
if os.IsNotExist(err) {
    return <b>stacktrace.PropagateWithCode(err, EcodeManifestNotFound, "")</b>
}
</pre>

#### stacktrace.NewMessageWithCode(code ErrorCode, msg string, vals ...interface{}) error

The error code mechanism can be useful by itself even where stack traces with
line numbers are not required. NewMessageWithCode returns an error that prints
just like `fmt.Errorf` with no line number, but including a code.

<pre>
ttl := req.URL.Query().Get("ttl")
if ttl == "" {
    return 0, <b>stacktrace.NewMessageWithCode(EcodeBadInput, "Missing ttl query parameter")</b>
}
</pre>

#### stacktrace.GetCode(err error) ErrorCode

GetCode extracts the error code from an error.

<pre>
for i := 0; i < attempts; i++ {
    err := Do()
    if <b>stacktrace.GetCode(err)</b> != EcodeTimeout {
        return err
    }
    // try a few more times
}
return stacktrace.NewError("timed out after %d attempts", attempts)
</pre>

GetCode returns the special value `stacktrace.NoCode` if `err` is nil or if
there is no error code attached to `err`.

## License

Stacktrace is released by Palantir Technologies, Inc. under the Apache 2.0
License. See the included [LICENSE](LICENSE) file for details.

## Contributing

We welcome contributions of backward-compatible changes to this library.

- Write your code
- Add tests for new functionality
- Run `go test` and verify that the tests pass
- Fill out the [Individual](https://github.com/palantir/stacktrace/blob/master/Palantir_Individual_Contributor_License_Agreement.pdf?raw=true) or [Corporate](https://github.com/palantir/stacktrace/blob/master/Palantir_Corporate_Contributor_License_Agreement.pdf?raw=true) Contributor License Agreement and send it to [opensource@palantir.com](mailto:opensource@palantir.com)
- Submit a pull request
