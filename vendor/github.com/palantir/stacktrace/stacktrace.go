// Copyright 2016 Palantir Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stacktrace

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/palantir/stacktrace/cleanpath"
)

/*
CleanPath function is applied to file paths before adding them to a stacktrace.
By default, it makes the path relative to the $GOPATH environment variable.

To remove some additional prefix like "github.com" from file paths in
stacktraces, use something like:

	stacktrace.CleanPath = func(path string) string {
		path = cleanpath.RemoveGoPath(path)
		path = strings.TrimPrefix(path, "github.com/")
		return path
	}
*/
var CleanPath = cleanpath.RemoveGoPath

/*
NewError is a drop-in replacement for fmt.Errorf that includes line number
information. The canonical call looks like this:

	if !IsOkay(arg) {
		return stacktrace.NewError("Expected %v to be okay", arg)
	}
*/
func NewError(msg string, vals ...interface{}) error {
	return create(nil, NoCode, msg, vals...)
}

/*
Propagate wraps an error to include line number information. The msg and vals
arguments work like the ones for fmt.Errorf.

The message passed to Propagate should describe the action that failed,
resulting in the cause. The canonical call looks like this:

	result, err := process(arg)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed to process %v", arg)
	}

To write the message, ask yourself "what does this call do?" What does
process(arg) do? It processes ${arg}, so the message is that we failed to
process ${arg}.

Pay attention that the message is not redundant with the one in err. If it is
not possible to add any useful contextual information beyond what is already
included in an error, msg can be an empty string:

	func Something() error {
		mutex.Lock()
		defer mutex.Unlock()

		err := reallySomething()
		return stacktrace.Propagate(err, "")
	}

If cause is nil, Propagate returns nil. This allows elision of some "if err !=
nil" checks.
*/
func Propagate(cause error, msg string, vals ...interface{}) error {
	if cause == nil {
		// Allow calling Propagate without checking whether there is error
		return nil
	}
	return create(cause, NoCode, msg, vals...)
}

/*
ErrorCode is a code that can be attached to an error as it is passed/propagated
up the stack.

There is no predefined set of error codes. You define the ones relevant to your
application:

	const (
		EcodeManifestNotFound = stacktrace.ErrorCode(iota)
		EcodeBadInput
		EcodeTimeout
	)

The one predefined error code is NoCode, which has a value of math.MaxUint16.
Avoid using that value as an error code.

An ordinary stacktrace.Propagate call preserves the error code of an error.
*/
type ErrorCode uint16

/*
NoCode is the error code of errors with no code explicitly attached.
*/
const NoCode ErrorCode = math.MaxUint16

/*
NewErrorWithCode is similar to NewError but also attaches an error code.
*/
func NewErrorWithCode(code ErrorCode, msg string, vals ...interface{}) error {
	return create(nil, code, msg, vals...)
}

/*
PropagateWithCode is similar to Propagate but also attaches an error code.

	_, err := os.Stat(manifestPath)
	if os.IsNotExist(err) {
		return stacktrace.PropagateWithCode(err, EcodeManifestNotFound, "")
	}
*/
func PropagateWithCode(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if cause == nil {
		// Allow calling PropagateWithCode without checking whether there is error
		return nil
	}
	return create(cause, code, msg, vals...)
}

/*
NewMessageWithCode returns an error that prints just like fmt.Errorf with no
line number, but including a code. The error code mechanism can be useful by
itself even where stack traces with line numbers are not warranted.

	ttl := req.URL.Query().Get("ttl")
	if ttl == "" {
		return 0, stacktrace.NewMessageWithCode(EcodeBadInput, "Missing ttl query parameter")
	}
*/
func NewMessageWithCode(code ErrorCode, msg string, vals ...interface{}) error {
	return &stacktrace{
		message: fmt.Sprintf(msg, vals...),
		code:    code,
	}
}

/*
GetCode extracts the error code from an error.

	for i := 0; i < attempts; i++ {
		err := Do()
		if stacktrace.GetCode(err) != EcodeTimeout {
			return err
		}
		// try a few more times
	}
	return stacktrace.NewError("timed out after %d attempts", attempts)

GetCode returns the special value stacktrace.NoCode if err is nil or if there is
no error code attached to err.
*/
func GetCode(err error) ErrorCode {
	if err, ok := err.(*stacktrace); ok {
		return err.code
	}
	return NoCode
}

type stacktrace struct {
	message  string
	cause    error
	code     ErrorCode
	file     string
	function string
	line     int
}

func create(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	// If no error code specified, inherit error code from the cause.
	if code == NoCode {
		code = GetCode(cause)
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, vals...),
		cause:   cause,
		code:    code,
	}

	// Caller of create is NewError or Propagate, so user's code is 2 up.
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	if CleanPath != nil {
		file = CleanPath(file)
	}
	err.file, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}
	err.function = shortFuncName(f)

	return err
}

/* "FuncName" or "Receiver.MethodName" */
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

func (st *stacktrace) Error() string {
	return fmt.Sprint(st)
}
