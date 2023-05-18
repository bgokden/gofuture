package gofuture

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Future type holds Result and state
type Future[T any] struct {
	Success          bool
	Done             bool
	Result           T
	InterfaceChannel <-chan T
	err              error
	errChannel       <-chan error
}

// Get return the result when available. This is a blocking call
func (f *Future[T]) Get() (T, error) {
	if f.Done {
		return f.Result, f.err
	}
	f.Result = <-f.InterfaceChannel
	f.err = <-f.errChannel
	f.Success = true
	f.Done = true
	return f.Result, f.err
}

// GetWithTimeout return the result until timeout.
func (f *Future[T]) GetWithTimeout(timeout time.Duration) (T, error) {
	if f.Done {
		return f.Result, f.err
	}
	timeoutChannel := time.After(timeout)
	select {
	case res := <-f.InterfaceChannel:
		f.Result = res
		f.err = <-f.errChannel
		f.Success = true
		f.Done = true
	case <-timeoutChannel:
		f.Result = reflect.Zero(reflect.TypeOf(f.Result)).Interface().(T)
		f.Done = true
		f.Success = false
		f.err = errors.New("timed out")
	}
	return f.Result, f.err
}

// FutureFunc creates a function that returns its response in future
func FutureFunc[T any](implem interface{}, args ...interface{}) *Future[T] {
	valIn := make([]reflect.Value, len(args), len(args))

	fnVal := reflect.ValueOf(implem)

	for idx, elt := range args {
		valIn[idx] = reflect.ValueOf(elt)
	}
	interfaceChannel := make(chan T, 1)
	errChannel := make(chan error, 1)
	go func() {
		defer func() {
			// handle the panic here
			if r := recover(); r != nil {
				errChannel <- fmt.Errorf("panic: %v", r)
			}
			close(interfaceChannel)
			close(errChannel)
		}()
		res := fnVal.Call(valIn)
		// Up to two return values are supported
		if len(res) > 1 && !res[1].IsNil() {
			// handle err
			errChannel <- res[1].Interface().(error)
		} else {
			// handle err
			errChannel <- nil
		}
		interfaceChannel <- res[0].Interface().(T)
	}()

	return &Future[T]{
		Success:          false,
		Done:             false,
		Result:           reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T),
		InterfaceChannel: interfaceChannel,
		errChannel:       errChannel,
	}
}
