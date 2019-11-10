package gofuture

import (
	"reflect"
	"time"
)

// Future type holds Result and state
type Future struct {
	Success          bool
	Done             bool
	Result           interface{}
	InterfaceChannel <-chan interface{}
}

// Get return the result when available. This is a blocking call
func (f *Future) Get() interface{} {
	if f.Done {
		return f.Result
	}
	f.Result = <-f.InterfaceChannel
	f.Success = true
	f.Done = true
	return f.Result
}

// GetWithTimeout return the result until timeout.
func (f *Future) GetWithTimeout(timeout time.Duration) interface{} {
	if f.Done {
		return f.Result
	}
	timeoutChannel := time.After(timeout)
	select {
	case res := <-f.InterfaceChannel:
		f.Result = res
		f.Success = true
		f.Done = true
	case <-timeoutChannel:
		f.Result = nil
		f.Done = true
		f.Success = false
	}
	return f.Result
}

// FutureFunc creates a function that returns its response in future
func FutureFunc(implem interface{}, args ...interface{}) *Future {
	valIn := make([]reflect.Value, len(args), len(args))

	fnVal := reflect.ValueOf(implem)

	for idx, elt := range args {
		valIn[idx] = reflect.ValueOf(elt)
	}
	interfaceChannel := make(chan interface{}, 1)

	go func() {
		res := fnVal.Call(valIn)
		// Only one result is supported
		interfaceChannel <- res[0].Interface()
	}()

	return &Future{
		Success:          false,
		Done:             false,
		Result:           nil,
		InterfaceChannel: interfaceChannel,
	}
}
