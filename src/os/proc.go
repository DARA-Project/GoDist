// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Process etc.

package os

import (
	"dara"
	"runtime"
	"syscall"
)

// Args hold the command-line arguments, starting with the program name.
var Args []string

func init() {
	if runtime.GOOS == "windows" {
		// Initialized in exec_windows.go.
		return
	}
	Args = runtime_args()
}

func runtime_args() []string // in package runtime

// Getuid returns the numeric user id of the caller.
//
// On Windows, it returns -1.
func Getuid() int {
	id := syscall.Getuid()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[GETUID]")
		retInfo := dara.GeneralType{Type:dara.INTEGER, Integer: id}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETUID, 0, 1, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETUID, syscallInfo)
	}
	return id
}

// Geteuid returns the numeric effective user id of the caller.
//
// On Windows, it returns -1.
func Geteuid() int {
	id := syscall.Geteuid()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[GETEUID]")
		retInfo := dara.GeneralType{Type:dara.INTEGER, Integer: id}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETEUID, 0, 1, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETEUID, syscallInfo)
	}
	return id
}

// Getgid returns the numeric group id of the caller.
//
// On Windows, it returns -1.
func Getgid() int {
	id := syscall.Getgid()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[GETGID]")
		retInfo := dara.GeneralType{Type:dara.INTEGER, Integer: id}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETGID, 0, 1, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETGID, syscallInfo)
	}
	return id
}

// Getegid returns the numeric effective group id of the caller.
//
// On Windows, it returns -1.
func Getegid() int {
	id := syscall.Getegid()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[GETEGID]")
		retInfo := dara.GeneralType{Type:dara.INTEGER, Integer: id}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETEGID, 0, 1, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETEGID, syscallInfo)
	}
	return id
}

// Getgroups returns a list of the numeric ids of groups that the caller belongs to.
//
// On Windows, it returns syscall.EWINDOWS. See the os/user package
// for a possible alternative.
func Getgroups() ([]int, error) {
	gids, e := syscall.Getgroups()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[GETGROUPS]")
		retInfo1 := dara.GeneralType{Type: dara.ARRAY, Integer: len(gids)}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETGROUPS, 0, 2, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETGROUPS, syscallInfo)
	}
	return gids, NewSyscallError("getgroups", e)
}

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.
func Exit(code int) {
	if runtime.Is_dara_profiling_on() {
		print("[EXIT] : ")
		println(code)
		argInfo := dara.GeneralType{Type: dara.INTEGER, Integer: code}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_EXIT, 1, 0, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_EXIT, syscallInfo)
	}
	if code == 0 {
		// Give race detector a chance to fail the program.
		// Racy programs do not have the right to finish successfully.
		runtime_beforeExit()
	}
	syscall.Exit(code)
}

func runtime_beforeExit() // implemented in runtime
