package gofuture_test

import (
	"fmt"
	"testing"
	"time"

	gofuture "github.com/bgokden/gofuture"
	"github.com/stretchr/testify/assert"
)

func TestFutureFunc(t *testing.T) {
	printTime(t)
	x := 10
	response := gofuture.FutureFunc(10*time.Second, func() int {
		printTime(t)
		time.Sleep(5 * time.Second)
		fmt.Printf("x = %v\n", x)
		return x * 10
	})
	printTime(t)

	result := response().Result

	t.Logf("Result: %v\n", result)
	printTime(t)

	assert.Equal(t, int(100), result)
}

func printTime(t *testing.T) {
	t.Logf("time: %v\n", time.Now().Unix())
}
