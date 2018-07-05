// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package os

import (
	"dara"
	"errors"
	"runtime"
	"syscall"
	"time"
)

func (p *Process) wait() (ps *ProcessState, err error) {
	if p.Pid == -1 {
		return nil, syscall.EINVAL
	}

	// If we can block until Wait4 will succeed immediately, do so.
	ready, err := p.blockUntilWaitable()
	if err != nil {
		return nil, err
	}
	if ready {
		// Mark the process done now, before the call to Wait4,
		// so that Process.signal will not send a signal.
		p.setDone()
		// Acquire a write lock on sigMu to wait for any
		// active call to the signal method to complete.
		p.sigMu.Lock()
		p.sigMu.Unlock()
	}

	var status syscall.WaitStatus
	var rusage syscall.Rusage
	pid1, e := syscall.Wait4(p.Pid, &status, 0, &rusage)
	if e != nil {
		// DARA Instrumentation
		if runtime.Is_dara_profiling_on() {
			print("[WAIT] : ")
			println(p.Pid)
			argInfo := dara.GeneralType{Type: dara.PROCESS, Integer: p.Pid}
			retInfo1 := dara.GeneralType{Type: dara.POINTER, Unsupported: dara.UNSUPPORTEDVAL}
			retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
			syscallInfo := dara.GeneralSyscall{dara.DSYS_WAIT4, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
			runtime.Report_Syscall_To_Scheduler(dara.DSYS_WAIT4, syscallInfo)
		}
		return nil, NewSyscallError("wait", e)
	}
	if pid1 != 0 {
		p.setDone()
	}
	ps = &ProcessState{
		pid:    pid1,
		status: status,
		rusage: &rusage,
	}
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[WAIT] : ")
		println(p.Pid)
		argInfo := dara.GeneralType{Type: dara.PROCESS, Integer: p.Pid}
		retInfo1 := dara.GeneralType{Type: dara.POINTER, Unsupported: dara.UNSUPPORTEDVAL}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_WAIT4, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_WAIT4, syscallInfo)
	}
	return ps, nil
}

var errFinished = errors.New("os: process already finished")

func (p *Process) signal(sig Signal) error {
	if p.Pid == -1 {
		return errors.New("os: process already released")
	}
	if p.Pid == 0 {
		return errors.New("os: process not initialized")
	}
	p.sigMu.RLock()
	defer p.sigMu.RUnlock()
	if p.done() {
		return errFinished
	}
	s, ok := sig.(syscall.Signal)
	if !ok {
		return errors.New("os: unsupported signal type")
	}
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[KILL] : ")
		print(p.Pid)
		print(" ")
		println(s.String())
		argInfo1 := dara.GeneralType{Type: dara.PROCESS, Integer: p.Pid}
		argInfo2 := dara.GeneralType{Type: dara.SIGNAL, String: s.String()}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo :=  dara.GeneralSyscall{dara.DSYS_KILL, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_KILL, syscallInfo)
	}
	if e := syscall.Kill(p.Pid, s); e != nil {
		if e == syscall.ESRCH {
			return errFinished
		}
		return e
	}
	return nil
}

func (p *Process) release() error {
	// NOOP for unix.
	p.Pid = -1
	// no need for a finalizer anymore
	runtime.SetFinalizer(p, nil)
	return nil
}

func findProcess(pid int) (p *Process, err error) {
	// NOOP for unix.
	return newProcess(pid, 0), nil
}

func (p *ProcessState) userTime() time.Duration {
	return time.Duration(p.rusage.Utime.Nano()) * time.Nanosecond
}

func (p *ProcessState) systemTime() time.Duration {
	return time.Duration(p.rusage.Stime.Nano()) * time.Nanosecond
}
