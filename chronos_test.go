package chronos

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// go test -v -run=XXX -bench=$ --race
// goos: windows
// goarch: amd64
// pkg: github.com/kataras/chronos
// BenchmarkChronos-8             1        4000087000 ns/op            1928 B/op         13 allocs/op
func BenchmarkChronos(b *testing.B) {
	var (
		c = New(5, 4*time.Second)
		i uint32
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i = 1; i <= c.Max+1; i++ {
		// it will add the first 5 at-hoc and wait 4 secs to add the +1.
		<-c.Acquire()
	}
}

// go test -v --race
const inspect = true

func log(format string, a ...interface{}) {
	if inspect {
		fmt.Printf(format, a...)
	}
}

const second = int64(time.Second)

type testCase struct {
	circle     uint64
	position   uint32
	sleepAfter int64
}

func TestChronos(t *testing.T) {
	var (
		max uint32 = 3
		per        = time.Second * 2
	)

	c := &C{
		Max: max,
		Per: per.Nanoseconds(),
	}

	tests := []testCase{
		{0, 1, 0},                     // 1
		{0, 2, 0},                     // 2
		{0, 3, 0},                     // 3
		{1, 1, 0},                     // 4
		{1, 2, per.Nanoseconds()},     // 5
		{2, 1, per.Nanoseconds() / 2}, // 6
		{2, 2, 0},                     // 7
		{2, 3, 0},                     // 8
		{3, 1, 0},                     // 9
		{3, 2, 0},                     // 10
	}

	do := func(c *C, i int, tt testCase) {
		<-c.Acquire()
		if expected, got := tt.circle, c.Circle(); expected != got {
			t.Fatalf("[%d] expected circle to be %d but got %d", i, expected, got)
		}
		if expected, got := tt.position, c.getCurrentLength(); expected != got {
			t.Fatalf("[%d] expected position to be %d but got %d", i, expected, got)
		}

		log("[%d] Acquire\n", i+1)

		if st := tt.sleepAfter; st > 0 {
			time.Sleep(time.Duration(st))
		}
	}

	for i, tt := range tests {
		do(c, i, tt)
	}
}

// TestChronosSchedDuration tests if the scheduler works as expected
// for items that exceeds the maximum given size.
func TestChronosSchedDuration(t *testing.T) {
	var (
		c = New(5, 4*time.Second)
		i uint32
	)

	for i = 1; i <= c.Max+1; i++ {
		testChronosSchedDuration(t, c, i)
	}

	log("---------ACQUIRE FROM %d DIFFERENT GOROUTINES---------\n", c.Max+1)
	c = New(5, 4*time.Second)
	var wg sync.WaitGroup
	for i = 1; i <= c.Max+1; i++ {
		wg.Add(1)
		go func(i uint32) {
			defer wg.Done()
			testChronosSchedDuration(t, c, i)
		}(i)
	}

	wg.Wait()
	// if test has max < per
	// then comment the wg things and uncomment this:
	// time.Sleep(time.Duration(int64(c.Max) * c.Per))
}

func testChronosSchedDuration(t *testing.T, c *C, i uint32) {
	now := time.Now()
	beforeFireLen := c.getCurrentLength()
	// that works as well as expected
	// if i%2 == 0 {
	// 	time.Sleep(time.Duration(c.per))
	// }
	<-c.Acquire()
	since := time.Since(now)

	log("[%d] Fired\n", i)

	// on differnet go routines this "i" is not reliable --> if i == c.max+1,
	// to check if this is the last item
	// we will use the current size.
	if beforeFireLen == c.Max {
		// the last should take c.per because it exceed the max size.
		// we will calculate all these in fixed(int) seconds because the execution
		// time may be delayed as you can imagine.
		if expected, got := int(c.Per/second), int(since.Seconds()); expected != got {
			t.Fatalf("expected to fire after %ds but fired after %ds\n", expected, got)
		}

		if got := c.getCurrentLength(); got != 1 {
			t.Fatalf("expected the length to be one, now that we are on different circle(%d) and that was the first acquire, but got %d", c.Circle(), got)
		}
	}
}
