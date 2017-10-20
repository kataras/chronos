package main

import (
	"fmt"
	"time"

	"github.com/kataras/chronos"
)

func main() {
	c := &chronos.C{
		Max: 3,
		Per: 5 * int64(time.Second),
		// or 5 * (time.Second.Nanoseconds())
		// or (5 * time.Second).Nanoseconds()
	}

	do(c) // [1] circle = 0 length = 1
	do(c) // [2] circle = 0 length = 2
	do(c) // [3] circle = 0 length = 3
	fmt.Println("waiting 5 seconds because of the acquire (if we had delay inside the function it wouldn't touch the chronos' time)")
	do(c) // [4] circle = 1 length = 1
	do(c) // [5] circle = 1 length = 2
	fmt.Println("wait again but this time because of the time.Sleep, the circle should change and the func should be fired without chronos' sched time")
	time.Sleep(time.Duration(c.Per))
	do(c) // [6] circle = 2 length = 1
	// fmt.Println("circle changed, do(c) called 3 times at 5 seconds")
	do(c) // [7] circle = 2 length = 2
	do(c) // [8] circle = 2 length = 3
	fmt.Println("wait for 5 seconds because of acquire, circle should be changed because of item exceed and duration passed")
	do(c) // [9] circle = 3 length = 1
	do(c) // [10] circle = 3 length = 2
}

var i int

func do(c *chronos.C) {
	i++
	<-c.Acquire()
	fmt.Printf("[%d] do inside circle: %d\n", i, c.Circle())
}
