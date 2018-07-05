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
		argInfo := dara.GeneralType{Type: dara.STRING, String: f.name}
		retInfo1 := dara.GeneralType{Type: dara.FILEINFO, Unsupported: dara.UNSUPPORTEDVAL}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FSTAT, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FSTAT, syscallInfo)
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
		argInfo := dara.GeneralType{Type: dara.STRING, String: name}
		retInfo1 := dara.GeneralType{Type: dara.FILEINFO, Unsupported: dara.UNSUPPORTEDVAL}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_STAT, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_STAT, syscallInfo)
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
		argInfo := dara.GeneralType{Type: dara.STRING, String: name}
		retInfo1 := dara.GeneralType{Type: dara.FILEINFO, Unsupported: dara.UNSUPPORTEDVAL}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_LSTAT, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_LSTAT, syscallInfo)
	}
	err := syscall.Lstat(name, &fs.sys)
	if err != nil {
		return nil, &PathError{"lstat", name, err}
	}
	fillFileStatFromSys(&fs, name)
	return &fs, nil
}
