package function

import (
	"fmt"
	"time"
)

func ExampleFunction_Call() {
	// 3  maximum functions are allowed to be executed
	// every 5 seconds.
	var maxCalls = 3
	per := 5 * time.Second
	f := New(uint32(maxCalls), per)

	var lastExecutedAt time.Time

	// we sleep 1 second per action but this doesn't matter
	// 3 calls will be able to executed at 5 seconds.
	action := func(message string) {
		// time.Sleep(1 * time.Second)
		fmt.Println(message)
	}

	time.AfterFunc(per-166*time.Millisecond, func() {
		fmt.Println("waiting for message-4 to be fired after 5 seconds the last one (or immediately if the limit time passed)...")
	})

	for i := 1; i <= maxCalls+1; i++ {
		lastExecutedAt = time.Now()
		message := fmt.Sprintf("message-%d", i)
		// output, err := f.Call(action, message)
		// <-output
		// or
		// output := <-f.MustCall(action, message)
		// or
		Wait(f.Call(action, message))
		if i == maxCalls+1 {
			since := time.Since(lastExecutedAt)
			// the last one will be executed after 5 seconds from the last executed.
			fmt.Printf("finish-after-%.0f-seconds", since.Seconds())
		}
	}

	// Output:
	// message-1
	// message-2
	// message-3
	// waiting for message-4 to be fired after 5 seconds the last one (or immediately if the limit time passed)...
	// message-4
	// finish-after-5-seconds
}
