package function

import (
	"reflect"
)

var interfaceOfError = reflect.TypeOf((*error)(nil)).Elem()

// LookupError tries to find an error value
// from the "out" values, if exists then it will return that
// otherwise nil.
func LookupError(out []reflect.Value) error {
	// we have the output parameters, check if something of those
	// is error type, and fill the channel with the first error type.
	for _, a := range out {
		if a.Type().Implements(interfaceOfError) {
			if !a.IsNil() && a.CanInterface() {
				// check again
				if err, ok := a.Interface().(error); ok {
					return err
				}
			}
		}
	}

	return nil
}

func AsError(ch <-chan []reflect.Value, err error) error {
	out := <-ch
	if err != nil {
		return err // couldn't even executed for some reason.
	}
	return LookupError(out)
}

type Result struct {
	Err error
	Out []reflect.Value
}

// OutErr returns the error of the actionFunc's Output, if any.
func (r Result) OutErr() error {
	return LookupError(r.Out)
}

// AsResult blocks until the "ch" to be ready
// it's used to wrap the `Call` function
// and return a `Result` struct which may be used
// to get a struct which contains some utils
// for the shake of simplicy,
// i.e: the error (if any) from the "Call's actionFunc" result.
func AsResult(ch <-chan []reflect.Value, err error) Result {
	out := <-ch

	return Result{
		Err: err,
		Out: out,
	}
}

func Wait(ch <-chan []reflect.Value, err error) ([]reflect.Value, error) {
	out := <-ch
	return out, err
}

/*
func showMeTheUsage() {
	fc := NewFuncCaller(nil)

	outAndDone, err := fn.Call(nil,nil)
	_ = <- outAndDone
	_ = outAndDone
	_ = err

	err = AsError(fn.Call(nil, nil))
	_ = err

	result := AsResult(fn.Call(nil, nil))
	_ = result

	outAndDone, err = Wait(fn.Call(nil,nil))
	_ = outAndDone
	_ = err
}
*/
