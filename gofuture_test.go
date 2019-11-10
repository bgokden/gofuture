package gofuture_test

import (
	"fmt"
	"testing"
	"time"

	gofuture "github.com/bgokden/gofuture"
	"github.com/stretchr/testify/assert"
)

func TestFutureFunc(t *testing.T) {
	x := 10
	var elapsed time.Duration
	start := time.Now()
	future := gofuture.FutureFunc(func() int {
		printTime(t)
		time.Sleep(5 * time.Second)
		fmt.Printf("x = %v\n", x)
		return x * 10
	})
	elapsed = time.Since(start)

	t.Logf("it took %s", elapsed)
	assert.Less(t, elapsed.Milliseconds(), (1 * time.Second).Milliseconds())

	result := future.Get()
	elapsed = time.Since(start)
	assert.Less(t, (5 * time.Second).Milliseconds(), elapsed.Milliseconds())
	assert.Equal(t, int(100), result)
	t.Logf("Result: %v\n", future.Get())

	// This assert tests calling result second time doesn't cause any problems
	assert.Equal(t, int(100), future.Get())
}

func TestFutureFuncTimeOut(t *testing.T) {
	x := 10
	var elapsed time.Duration
	start := time.Now()
	future := gofuture.FutureFunc(func() int {
		printTime(t)
		time.Sleep(5 * time.Second)
		fmt.Printf("x = %v\n", x)
		return x * 10
	})
	elapsed = time.Since(start)

	t.Logf("it took %s", elapsed)
	assert.Less(t, elapsed.Milliseconds(), (1 * time.Second).Milliseconds())

	result := future.GetWithTimeout(3 * time.Second)
	elapsed = time.Since(start)
	assert.Equal(t, nil, result)
	assert.Less(t, elapsed.Milliseconds(), (4 * time.Second).Milliseconds())
	t.Logf("Result: %v\n", result)

	assert.Equal(t, nil, future.Get())
}

func printTime(t *testing.T) {
	t.Logf("time: %v\n", time.Now().Unix())
}
