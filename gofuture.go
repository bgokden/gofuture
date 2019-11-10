package gofuture

import (
	"fmt"
	"reflect"
	"time"
)

type Future struct {
	Success          bool
	Done             bool
	Result           interface{}
	InterfaceChannel <-chan interface{}
	TimeoutChannel   <-chan time.Time
}

func (f *Future) Get() interface{} {
	if f.Done {
		return f.Result
	}
	select {
	case res := <-f.InterfaceChannel:
		fmt.Println(res)
		f.Result = res
		f.Success = true
		f.Done = true
		return res
	case <-f.TimeoutChannel:
		f.Result = nil
		f.Done = true
		f.Success = false
		return nil
	}
}

// FutureFunc creates a function that returns its response in future
func FutureFunc(timeout time.Duration, implem interface{}, args ...interface{}) *Future {
	valIn := make([]reflect.Value, len(args), len(args))

	fnVal := reflect.ValueOf(implem)

	for idx, elt := range args {
		valIn[idx] = reflect.ValueOf(elt)
	}
	interfaceChannel := make(chan interface{}, 1)
	timeoutChannel := time.After(timeout)

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
		TimeoutChannel:   timeoutChannel,
	}
}
