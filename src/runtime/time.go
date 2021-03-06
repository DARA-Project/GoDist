// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Time-related runtime and pieces of package time.

package runtime

import (
	"runtime/internal/sys"
	"unsafe"
    "dara"
)

// Package time knows the layout of this structure.
// If this struct changes, adjust ../time/sleep.go:/runtimeTimer.
// For GOOS=nacl, package syscall knows the layout of this structure.
// If this struct changes, adjust ../syscall/net_nacl.go:/runtimeTimer.
type timer struct {
	tb *timersBucket // the bucket the timer lives in
	i  int           // heap index

	// Timer wakes up at when, and then at when+period, ... (period > 0 only)
	// each time calling f(arg, now) in the timer goroutine, so f must be
	// a well-behaved function and not block.
	when   int64
	period int64
	f      func(interface{}, uintptr)
	arg    interface{}
	seq    uintptr
}

// timersLen is the length of timers array.
//
// Ideally, this would be set to GOMAXPROCS, but that would require
// dynamic reallocation
//
// The current value is a compromise between memory usage and performance
// that should cover the majority of GOMAXPROCS values used in the wild.
const timersLen = 64

// timers contains "per-P" timer heaps.
//
// Timers are queued into timersBucket associated with the current P,
// so each P may work with its own timers independently of other P instances.
//
// Each timersBucket may be associated with multiple P
// if GOMAXPROCS > timersLen.
var timers [timersLen]struct {
	timersBucket

	// The padding should eliminate false sharing
	// between timersBucket values.
	pad [sys.CacheLineSize - unsafe.Sizeof(timersBucket{})%sys.CacheLineSize]byte
}

func (t *timer) assignBucket() *timersBucket {
	id := uint8(getg().m.p.ptr().id) % timersLen
	t.tb = &timers[id].timersBucket
	return t.tb
}

//go:notinheap
type timersBucket struct {
	lock         mutex
	gp           *g
	created      bool
	sleeping     bool
	rescheduling bool
	sleepUntil   int64
	waitnote     note
	t            []*timer
}

// nacl fake time support - time in nanoseconds since 1970
var faketime int64

// Package time APIs.
// Godoc uses the comments in package time, not these.

// time.now is implemented in assembly.

// timeSleep puts the current goroutine to sleep for at least ns nanoseconds.
//go:linkname timeSleep time.Sleep
func timeSleep(ns int64) {
	// DARA: Don't report the call if the time is less than 0....Basically is a no-op function call.
	if ns <= 0 {
		return
	}
	gp := getg()
	//Dara injection
    if Is_dara_profiling_on() {
        Dara_Debug_Print(func() { println("[SLEEP]", ns) })
		argInfo := dara.GeneralType{Type: dara.INTEGER64, Integer64: ns}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER64, Integer64: gp.goid}
        syscallInfo := dara.GeneralSyscall{dara.DSYS_SLEEP, 2, 0, [10]dara.GeneralType{argInfo, argInfo2}, [10]dara.GeneralType{}}
        Report_Syscall_To_Scheduler(dara.DSYS_SLEEP, syscallInfo)
    }
	t := gp.timer
	if t == nil {
		t = new(timer)
		gp.timer = t
	}
	*t = timer{}
    //Dara injection
    if Replay || Explore {
        dprint(dara.INFO, func () {println("[GoRoutine]timeSleep : Goroutine here for nap time")})
        //Don't install the timer in replay but obtain the lock :)
        tb := t.assignBucket()
        lock(&tb.lock)
        goparkunlock(&tb.lock, "sleep", traceEvGoSleep, 2)
        return
    }
    if Replay || Explore{
        dprint(dara.WARN, func () {println("[GoRuntime]timeSleep : We shouldn't be here wtf")})
    }
	t.when = nanotime() + ns
	t.f = goroutineReady
	t.arg = gp
	tb := t.assignBucket()
	lock(&tb.lock)
	tb.addtimerLocked(t)
	goparkunlock(&tb.lock, "sleep", traceEvGoSleep, 2)
}

// startTimer adds t to the timer heap.
//go:linkname startTimer time.startTimer
func startTimer(t *timer) {
	if raceenabled {
		racerelease(unsafe.Pointer(t))
	}
	addtimer(t)
}

// stopTimer removes t from the timer heap if it is there.
// It returns true if t was removed, false if t wasn't even there.
//go:linkname stopTimer time.stopTimer
func stopTimer(t *timer) bool {
	return deltimer(t)
}

// Go runtime.

// Ready the goroutine arg.
func goroutineReady(arg interface{}, seq uintptr) {
    dprint(dara.INFO, func() {println("[GoRoutine]goroutineReady : Goroutine that woke up :", arg.(*g).goid)})
	goready(arg.(*g), 0)
}

func addtimer(t *timer) {
	if DaraInitialised {
		// Add the timer to a list of timers which is exposed
		// to the global scheduler and have it choose firing off the timer as one of its actions.
		TimerCount += 1
		TimerInfo[TimerCount] = t
		LogTimerEvent(t)
		println(t.when)
	}
	if Explore {
		// Only prevent a timer from getting installed if this is exploration!
		return
	}
	if Record {
		// Fuck timers for now
		return
	}
	tb := t.assignBucket()
	lock(&tb.lock)
	tb.addtimerLocked(t)
	unlock(&tb.lock)
}

// Add a timer to the heap and start or kick timerproc if the new timer is
// earlier than any of the others.
// Timers are locked.
func (tb *timersBucket) addtimerLocked(t *timer) {
	// when must never be negative; otherwise timerproc will overflow
	// during its delta calculation and never expire other runtime timers.
	if t.when < 0 {
		t.when = 1<<63 - 1
	}
	t.i = len(tb.t)
	tb.t = append(tb.t, t)
	siftupTimer(tb.t, t.i)
	if t.i == 0 {
		// siftup moved to top: new earliest deadline.
		if tb.sleeping {
			tb.sleeping = false
			notewakeup(&tb.waitnote)
		}
		if tb.rescheduling {
			tb.rescheduling = false
			goready(tb.gp, 0)
		}
	}
	if !tb.created {
		tb.created = true
		go timerproc(tb)
	}
}

// Delete timer t from the heap.
// Do not need to update the timerproc: if it wakes up early, no big deal.
func deltimer(t *timer) bool {
	if t.tb == nil {
		// t.tb can be nil if the user created a timer
		// directly, without invoking startTimer e.g
		//    time.Ticker{C: c}
		// In this case, return early without any deletion.
		// See Issue 21874.
		return false
	}

	tb := t.tb

	lock(&tb.lock)
	// t may not be registered anymore and may have
	// a bogus i (typically 0, if generated by Go).
	// Verify it before proceeding.
	i := t.i
	last := len(tb.t) - 1
	if i < 0 || i > last || tb.t[i] != t {
		unlock(&tb.lock)
		return false
	}
	if i != last {
		tb.t[i] = tb.t[last]
		tb.t[i].i = i
	}
	tb.t[last] = nil
	tb.t = tb.t[:last]
	if i != last {
		siftupTimer(tb.t, i)
		siftdownTimer(tb.t, i)
	}
	unlock(&tb.lock)
	return true
}

// Timerproc runs the time-driven events.
// It sleeps until the next event in the tb heap.
// If addtimer inserts a new earlier event, it wakes timerproc early.
func timerproc(tb *timersBucket) {
	tb.gp = getg()
	for {
		lock(&tb.lock)
		tb.sleeping = false
		now := nanotime()
		delta := int64(-1)
		for {
			if len(tb.t) == 0 {
				delta = -1
				break
			}
			t := tb.t[0]
			delta = t.when - now
			if delta > 0 {
				break
			}
			if t.period > 0 {
				// leave in heap but adjust next time to fire
				t.when += t.period * (1 + -delta/t.period)
				siftdownTimer(tb.t, 0)
			} else {
				// remove from heap
				last := len(tb.t) - 1
				if last > 0 {
					tb.t[0] = tb.t[last]
					tb.t[0].i = 0
				}
				tb.t[last] = nil
				tb.t = tb.t[:last]
				if last > 0 {
					siftdownTimer(tb.t, 0)
				}
				t.i = -1 // mark as removed
			}
			f := t.f
			arg := t.arg
			seq := t.seq
			unlock(&tb.lock)
			if raceenabled {
				raceacquire(unsafe.Pointer(t))
			}
			f(arg, seq)
			lock(&tb.lock)
		}
		if delta < 0 || faketime > 0 {
			// No timers left - put goroutine to sleep.
			tb.rescheduling = true
			goparkunlock(&tb.lock, "timer goroutine (idle)", traceEvGoBlock, 1)
			continue
		}
		// At least one timer pending. Sleep until then.
		tb.sleeping = true
		tb.sleepUntil = now + delta
		noteclear(&tb.waitnote)
		unlock(&tb.lock)
		notetsleepg(&tb.waitnote, delta)
	}
}

func timejump() *g {
	if faketime == 0 {
		return nil
	}

	for i := range timers {
		lock(&timers[i].lock)
	}
	gp := timejumpLocked()
	for i := range timers {
		unlock(&timers[i].lock)
	}

	return gp
}

func timejumpLocked() *g {
	// Determine a timer bucket with minimum when.
	var minT *timer
	for i := range timers {
		tb := &timers[i]
		if !tb.created || len(tb.t) == 0 {
			continue
		}
		t := tb.t[0]
		if minT == nil || t.when < minT.when {
			minT = t
		}
	}
	if minT == nil || minT.when <= faketime {
		return nil
	}

	faketime = minT.when
	tb := minT.tb
	if !tb.rescheduling {
		return nil
	}
	tb.rescheduling = false
	return tb.gp
}

func timeSleepUntil() int64 {
	next := int64(1<<63 - 1)

	// Determine minimum sleepUntil across all the timer buckets.
	//
	// The function can not return a precise answer,
	// as another timer may pop in as soon as timers have been unlocked.
	// So lock the timers one by one instead of all at once.
	for i := range timers {
		tb := &timers[i]

		lock(&tb.lock)
		if tb.sleeping && tb.sleepUntil < next {
			next = tb.sleepUntil
		}
		unlock(&tb.lock)
	}

	return next
}

// Heap maintenance algorithms.

func siftupTimer(t []*timer, i int) {
	when := t[i].when
	tmp := t[i]
	for i > 0 {
		p := (i - 1) / 4 // parent
		if when >= t[p].when {
			break
		}
		t[i] = t[p]
		t[i].i = i
		i = p
	}
	if tmp != t[i] {
		t[i] = tmp
		t[i].i = i
	}
}

func siftdownTimer(t []*timer, i int) {
	n := len(t)
	when := t[i].when
	tmp := t[i]
	for {
		c := i*4 + 1 // left child
		c3 := c + 2  // mid child
		if c >= n {
			break
		}
		w := t[c].when
		if c+1 < n && t[c+1].when < w {
			w = t[c+1].when
			c++
		}
		if c3 < n {
			w3 := t[c3].when
			if c3+1 < n && t[c3+1].when < w3 {
				w3 = t[c3+1].when
				c3++
			}
			if w3 < w {
				w = w3
				c = c3
			}
		}
		if w >= when {
			break
		}
		t[i] = t[c]
		t[i].i = i
		i = c
	}
	if tmp != t[i] {
		t[i] = tmp
		t[i].i = i
	}
}

// Entry points for net, time to call nanotime.

//go:linkname poll_runtimeNano internal/poll.runtimeNano
func poll_runtimeNano() int64 {
	return nanotime()
}

//go:linkname time_runtimeNano time.runtimeNano
func time_runtimeNano() int64 {
	return nanotime()
}

// Monotonic times are reported as offsets from startNano.
// We initialize startNano to nanotime() - 1 so that on systems where
// monotonic time resolution is fairly low (e.g. Windows 2008
// which appears to have a default resolution of 15ms),
// we avoid ever reporting a nanotime of 0.
// (Callers may want to use 0 as "time not set".)
var startNano int64 = nanotime() - 1
