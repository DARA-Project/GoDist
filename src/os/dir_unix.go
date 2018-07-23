// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package os

import (
	"dara"
	"io"
	"runtime"
	"syscall"
)

const (
	blockSize = 4096
)

func (f *File) readdir(n int) (fi []FileInfo, err error) {
	dirname := f.name
	if dirname == "" {
		dirname = "."
	}
	names, err := f.Readdirnames(n)
	fi = make([]FileInfo, 0, len(names))
	for _, filename := range names {
		fip, lerr := lstat(dirname + "/" + filename)
		if IsNotExist(lerr) {
			// File disappeared between readdir + stat.
			// Just treat it as if it didn't exist.
			continue
		}
		if lerr != nil {
			if runtime.Is_dara_profiling_on() {
				print("[Readdir] : ")
				print(f.file.name)
				print(" ")
				println(n)
				argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
				argInfo2 := dara.GeneralType{Type:dara.INTEGER, Integer:n}
				retInfo1 := dara.GeneralType{Type:dara.ARRAY, Integer: len(fi)}
				retInfo2 := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
				syscallInfo := dara.GeneralSyscall{dara.DSYS_READDIR, 2, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo1, retInfo2}}
				runtime.Report_Syscall_To_Scheduler(dara.DSYS_READDIR, syscallInfo)
			}
			return fi, lerr
		}
		fi = append(fi, fip)
	}
	if len(fi) == 0 && err == nil && n > 0 {
		// Per File.Readdir, the slice must be non-empty or err
		// must be non-nil if n > 0.
		err = io.EOF
	}
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[Readdir] : ")
		print(f.file.name)
		print(" ")
		println(n)
		argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type:dara.INTEGER, Integer:n}
		retInfo1 := dara.GeneralType{Type:dara.ARRAY, Integer: len(fi)}
		retInfo2 := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_READDIR, 2, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_READDIR, syscallInfo)
	}
	return fi, err
}

func (f *File) readdirnames(n int) (names []string, err error) {
	// If this file has no dirinfo, create one.
	if f.dirinfo == nil {
		f.dirinfo = new(dirInfo)
		// The buffer must be at least a block long.
		f.dirinfo.buf = make([]byte, blockSize)
	}
	d := f.dirinfo

	size := n
	if size <= 0 {
		size = 100
		n = -1
	}

	names = make([]string, 0, size) // Empty with room to grow.
	for n != 0 {
		// Refill the buffer if necessary
		if d.bufp >= d.nbuf {
			d.bufp = 0
			var errno error
			d.nbuf, errno = f.pfd.ReadDirent(d.buf)
			runtime.KeepAlive(f)
			if errno != nil {
				if runtime.Is_dara_profiling_on() {
					print("[Readdirnames] : ")
					print(f.file.name)
					print(" ")
					println(n)
					argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
					argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer:n}
					retInfo1 := dara.GeneralType{Type: dara.ARRAY, Integer: len(names)}
					retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
					syscallInfo := dara.GeneralSyscall{dara.DSYS_READDIRNAMES, 2, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo1, retInfo2}}
					runtime.Report_Syscall_To_Scheduler(dara.DSYS_READDIRNAMES, syscallInfo)
				}
				return names, wrapSyscallError("readdirent", errno)
			}
			if d.nbuf <= 0 {
				break // EOF
			}
		}

		// Drain the buffer
		var nb, nc int
		nb, nc, names = syscall.ParseDirent(d.buf[d.bufp:d.nbuf], n, names)
		d.bufp += nb
		n -= nc
	}
	if n >= 0 && len(names) == 0 {
		if runtime.Is_dara_profiling_on() {
			print("[Readdirnames] : ")
			print(f.file.name)
			print(" ")
			println(n)
			argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
			argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer:n}
			retInfo1 := dara.GeneralType{Type: dara.ARRAY, Integer: len(names)}
			retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
			syscallInfo := dara.GeneralSyscall{dara.DSYS_READDIRNAMES, 2, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo1, retInfo2}}
			runtime.Report_Syscall_To_Scheduler(dara.DSYS_READDIRNAMES, syscallInfo)
		}
		return names, io.EOF
	}
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[Readdirnames] : ")
		print(f.file.name)
		print(" ")
		println(n)
		argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer:n}
		retInfo1 := dara.GeneralType{Type: dara.ARRAY, Integer: len(names)}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_READDIRNAMES, 2, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_READDIRNAMES, syscallInfo)
	}
	return names, nil
}
