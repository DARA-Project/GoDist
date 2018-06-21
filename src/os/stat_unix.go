// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package os

import (
	"dara"
	"runtime"
	"syscall"
)

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.
func (f *File) Stat() (FileInfo, error) {
    // DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FSTAT] : ")
		println(f.file.name)
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FSTAT)
	}
	if f == nil {
		return nil, ErrInvalid
	}
	var fs fileStat
	err := f.pfd.Fstat(&fs.sys)
	if err != nil {
		return nil, &PathError{"stat", f.name, err}
	}
	fillFileStatFromSys(&fs, f.name)
	return &fs, nil
}

// statNolog stats a file with no test logging.
func statNolog(name string) (FileInfo, error) {
	var fs fileStat
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[STAT] : " + name)
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_STAT)
	}
	err := syscall.Stat(name, &fs.sys)
	if err != nil {
		return nil, &PathError{"stat", name, err}
	}
	fillFileStatFromSys(&fs, name)
	return &fs, nil
}

// lstatNolog lstats a file with no test logging.
func lstatNolog(name string) (FileInfo, error) {
	var fs fileStat
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[LSTAT] : " + name)
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_LSTAT)
	}
	err := syscall.Lstat(name, &fs.sys)
	if err != nil {
		return nil, &PathError{"lstat", name, err}
	}
	fillFileStatFromSys(&fs, name)
	return &fs, nil
}
