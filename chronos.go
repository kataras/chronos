package chronos

import (
	"sync"
	"sync/atomic"
	"time"
)

// C is the one and only structure of the chronos package.
// It keeps the information about the X max operations per Y time
// and the `Acquire` function which is being used for availability "searching".
//
// Use the `New` or just the `&C{ Max: 5, Per: 3 * time.Second}`.
type C struct { // C is actually a kinda of rate limiter :)
	Max uint32 // maximum operations
	Per int64  // per x time (in nanoseconds).

	// starting from zero,
	// it's fair because when starting is not a complete circle:)
	circle uint64

	mu        sync.RWMutex
	length    uint32 // starting from zero but it will be 1 at first Acquire.
	lastAdded int64  // starting from zero time.
}

// New initializes and returns a new C chronos.
// The first input argument is the the maximum operations
// and the second is the time duration which should passed
// between the calls.
//
// X "max" operations "per" Y time duration.
func New(max uint32, per time.Duration) *C {
	return &C{
		Max: max,
		Per: int64(per), // nanoseconds
	}
}

func (c *C) getLastAdded() int64 {
	return atomic.LoadInt64(&c.lastAdded)
}

func (c *C) setLastAdded(t int64) {
	atomic.StoreInt64(&c.lastAdded, t)
}

func (c *C) getCurrentLength() uint32 {
	return atomic.LoadUint32(&c.length)
}

func (c *C) increment(delta uint32) uint32 {
	return atomic.AddUint32(&c.length, delta)
}

func (c *C) resetLength() {
	atomic.StoreUint32(&c.length, 0)
}

func (c *C) drawCircle() {
	atomic.AddUint64(&c.circle, 1)
}

// Circle returns the current "circle".
// A circle is changed when a new group of operations
// are called or when the sched duration passed and `Acquire` is called.
func (c *C) Circle() uint64 {
	return atomic.LoadUint64(&c.circle)
}

var emptyStruct = struct{}{}

func (c *C) tryAcquire(ch chan struct{}, inc bool) {
	current := time.Now().UnixNano()
	lastAdded := c.getLastAdded()

	if lastAdded != 0 && current-lastAdded-c.Per > 0 {
		c.drawCircle()
		c.resetLength()
	}

	// if the current length is smaller than the max
	// then we don't have to check for anything else,
	// it's available.
	// Remember: length starts from 0 when max from 1.
	if c.getCurrentLength() < c.Max {
		if inc { // no need of ratio between actions, yet.
			c.increment(1)
		}

		c.setLastAdded(current)
		ch <- emptyStruct
		return
	}

	// else schedule that.
	sched := current - lastAdded
	if sched <= c.Per {
		sched = c.Per - sched
	}

	time.AfterFunc(time.Duration(sched), func() {
		c.mu.Lock()
		c.tryAcquire(ch, true)
		c.mu.Unlock()
	})
}

// Acquire is the only one function of the chronos core.
// It will block if the already called times are > than the given "max" operations.
func (c *C) Acquire() <-chan struct{} {
	ch := make(chan struct{}, 0)
	go c.tryAcquire(ch, true)
	return ch
}
