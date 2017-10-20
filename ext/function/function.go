package function

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kataras/chronos"
)

type Function struct {
	*chronos.C
}

func New(max uint32, per time.Duration) *Function {
	return &Function{C: chronos.New(max, per)}
}

var Panic = func(err error) {
	panic(err)
}

func (f *Function) MustCall(actionFunc interface{}, actionFuncInput ...interface{}) <-chan []reflect.Value {
	done, err := f.Call(actionFunc, actionFuncInput...)
	if err != nil {
		Panic(err)
	}
	return done
}

func isGoodFunc(fn reflect.Value) bool { return fn.IsValid() && fn.Kind() == reflect.Func }

func (f *Function) Call(actionFunc interface{}, actionFuncInput ...interface{}) (<-chan []reflect.Value, error) {
	ch := make(chan []reflect.Value, 0)

	fn := reflect.ValueOf(actionFunc)
	if !isGoodFunc(fn) {
		ch <- nil
		return ch, fmt.Errorf("invalid kind, the 'actionFunc' should be a non-nil function with or without input arguments")
	}

	in := make([]reflect.Value, len(actionFuncInput), len(actionFuncInput))
	for i, a := range actionFuncInput {
		in[i] = reflect.ValueOf(a)
	}

	go f.call(fn, in, ch)
	return ch, nil
}

func (f *Function) call(fn reflect.Value, in []reflect.Value, ch chan []reflect.Value) {
	if fn.IsNil() {
		ch <- nil
		// something happened and fn now is nil, maybe the container was collected?
		// if so, then don't fire an error and just skip it.
		return
	}

	<-f.C.Acquire()
	ch <- fn.Call(in)
}
