package gofuture

import (
	"fmt"
	"reflect"
	"time"
)

type Future struct {
	Success bool
	Result  interface{}
}

func FutureFunc(timeout time.Duration, implem interface{}, args ...interface{}) func() Future {
	valIn := make([]reflect.Value, len(args), len(args))

	fnVal := reflect.ValueOf(implem)

	for idx, elt := range args {
		valIn[idx] = reflect.ValueOf(elt)
	}
	c1 := make(chan interface{}, 1)
	timeoutCh := time.After(timeout)

	// Run your long running function in it's own goroutine and pass back it's
	// response into our channel.
	go func() {
		res := fnVal.Call(valIn)
		c1 <- res[0].Interface()
	}()

	// Listen on our channel AND a timeout channel - which ever happens first.
	return func() Future {
		// Listen on our channel AND a timeout channel - whichever happens first.
		select {
		case res := <-c1:
			fmt.Println(res)
			return Future{
				Success: true,
				Result:  res,
			}
		case <-timeoutCh:
			return Future{
				Success: false,
				Result:  nil,
			}
		}
	}
}
