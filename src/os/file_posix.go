// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris windows

package os

import (
	"dara"
	"runtime"
	"syscall"
	"time"
)

func sigpipe() // implemented in package runtime

// Readlink returns the destination of the named symbolic link.
// If there is an error, it will be of type *PathError.
func Readlink(name string) (string, error) {
	for len := 128; ; len *= 2 {
		b := make([]byte, len)
		n, e := fixCount(syscall.Readlink(fixLongPath(name), b))
		if e != nil {
			// DARA Instrumentation
			if runtime.Is_dara_profiling_on() {
				print("[READLINK] : ")
				println(name)
				argInfo := dara.GeneralType{Type: dara.STRING, String: name}
				retInfo1 := dara.GeneralType{Type: dara.STRING, String: ""}
				retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported : dara.UNSUPPORTEDVAL}
				syscallInfo := dara.GeneralSyscall{dara.DSYS_READLINK, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
				runtime.Report_Syscall_To_Scheduler(dara.DSYS_READLINK, syscallInfo)
			}
			return "", &PathError{"readlink", name, e}
		}
		if n < len {
			// DARA Instrumentation
			if runtime.Is_dara_profiling_on() {
				print("[READLINK] : ")
				println(name)
				argInfo := dara.GeneralType{Type: dara.STRING, String: name}
				retInfo1 := dara.GeneralType{Type: dara.STRING, String: string(b[0:n])}
				retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported : dara.UNSUPPORTEDVAL}
				syscallInfo := dara.GeneralSyscall{dara.DSYS_READLINK, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
				runtime.Report_Syscall_To_Scheduler(dara.DSYS_READLINK, syscallInfo)
			}
			return string(b[0:n]), nil
		}
	}
}

// syscallMode returns the syscall-specific mode bits from Go's portable mode bits.
func syscallMode(i FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	// No mapping for Go's ModeTemporary (plan9 only).
	return
}

// See docs in file.go:Chmod.
func chmod(name string, mode FileMode) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[CHMOD] : ")
		print(name)
		print(" ")
		println(mode)
		argInfo1 := dara.GeneralType{Type: dara.STRING, String : name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer: int(mode)}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_CHMOD, 1, 2, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_CHMOD, syscallInfo)
	}
	if e := syscall.Chmod(fixLongPath(name), syscallMode(mode)); e != nil {
		return &PathError{"chmod", name, e}
	}
	return nil
}

// See docs in file.go:(*File).Chmod.
func (f *File) chmod(mode FileMode) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FCHMOD] : ")
		print(f.file.name)
		print(" ")
		println(mode)
		argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer: int(mode)}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FCHMOD, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FCHMOD, syscallInfo)
	}
	if err := f.checkValid("chmod"); err != nil {
		return err
	}
	if e := f.pfd.Fchmod(syscallMode(mode)); e != nil {
		return f.wrapErr("chmod", e)
	}
	return nil
}

// Chown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link's target.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func Chown(name string, uid, gid int) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[CHOWN] : ")
		print(name)
		print(" ")
		print(uid)
		print(" ")
		println(gid)
		argInfo1 := dara.GeneralType{Type: dara.STRING, String: name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer: int(uid)}
		argInfo3 := dara.GeneralType{Type: dara.INTEGER, Integer: int(gid)}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_CHOWN, 3, 1, [10]dara.GeneralType{argInfo1, argInfo2, argInfo3}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_CHOWN, syscallInfo)
	}
	if e := syscall.Chown(name, uid, gid); e != nil {
		return &PathError{"chown", name, e}
	}
	return nil
}

// Lchown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link itself.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func Lchown(name string, uid, gid int) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[LCHOWN] : ")
		print(name)
		print(" ")
		print(uid)
		print(" ")
		println(gid)
		argInfo1 := dara.GeneralType{Type: dara.STRING, String: name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer: int(uid)}
		argInfo3 := dara.GeneralType{Type: dara.INTEGER, Integer: int(gid)}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_LCHOWN, 3, 1, [10]dara.GeneralType{argInfo1, argInfo2, argInfo3}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_LCHOWN, syscallInfo)
	}
	if e := syscall.Lchown(name, uid, gid); e != nil {
		return &PathError{"lchown", name, e}
	}
	return nil
}

// Chown changes the numeric uid and gid of the named file.
// If there is an error, it will be of type *PathError.
//
// On Windows, it always returns the syscall.EWINDOWS error, wrapped
// in *PathError.
func (f *File) Chown(uid, gid int) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FCHOWN] : ")
		print(f.file.name)
		print(" ")
		print(uid)
		print(" ")
		println(gid)
		argInfo1 := dara.GeneralType{Type: dara.STRING, String: f.name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER, Integer: int(uid)}
		argInfo3 := dara.GeneralType{Type: dara.INTEGER, Integer: int(gid)}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FCHOWN, 3, 1, [10]dara.GeneralType{argInfo1, argInfo2, argInfo3}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FCHOWN, syscallInfo)
	}
	if err := f.checkValid("chown"); err != nil {
		return err
	}
	if e := f.pfd.Fchown(uid, gid); e != nil {
		return f.wrapErr("chown", e)
	}
	return nil
}

// Truncate changes the size of the file.
// It does not change the I/O offset.
// If there is an error, it will be of type *PathError.
func (f *File) Truncate(size int64) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FTRUNCATE] : ")
		print(f.file.name)
		print(" ")
		print(size)
		argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type: dara.INTEGER64, Integer64: size}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FTRUNCATE, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FTRUNCATE, syscallInfo)
	}
	if err := f.checkValid("truncate"); err != nil {
		return err
	}
	if e := f.pfd.Ftruncate(size); e != nil {
		return f.wrapErr("truncate", e)
	}
	return nil
}

// Sync commits the current contents of the file to stable storage.
// Typically, this means flushing the file system's in-memory copy
// of recently written data to disk.
func (f *File) Sync() error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FSYNC] : ")
		println(f.file.name)
		argInfo1 := dara.GeneralType{Type: dara.FILE, String: f.name}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FSYNC, 1, 1, [10]dara.GeneralType{argInfo1}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FSYNC, syscallInfo)
	}
	if err := f.checkValid("sync"); err != nil {
		return err
	}
	if e := f.pfd.Fsync(); e != nil {
		return f.wrapErr("sync", e)
	}
	return nil
}

// Chtimes changes the access and modification times of the named
// file, similar to the Unix utime() or utimes() functions.
//
// The underlying filesystem may truncate or round the values to a
// less precise time unit.
// If there is an error, it will be of type *PathError.
func Chtimes(name string, atime time.Time, mtime time.Time) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[UTIMES] : ")
		print(name)
		print(" ")
		print(atime.String())
		print(" ")
		println(mtime.String())
		argInfo1 := dara.GeneralType{Type:dara.STRING, String: name}
		argInfo2 := dara.GeneralType{Type:dara.TIME, String: atime.String()}
		argInfo3 := dara.GeneralType{Type:dara.TIME, String: mtime.String()}
		retInfo := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_UTIMES, 3, 1, [10]dara.GeneralType{argInfo1, argInfo2, argInfo3}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_UTIMES, syscallInfo)
	}
	var utimes [2]syscall.Timespec
	utimes[0] = syscall.NsecToTimespec(atime.UnixNano())
	utimes[1] = syscall.NsecToTimespec(mtime.UnixNano())
	if e := syscall.UtimesNano(fixLongPath(name), utimes[0:]); e != nil {
		return &PathError{"chtimes", name, e}
	}
	return nil
}

// Chdir changes the current working directory to the file,
// which must be a directory.
// If there is an error, it will be of type *PathError.
func (f *File) Chdir() error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[FCHDIR] : ")
		println(f.file.name)
		argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
		retInfo := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_FCHDIR, 1, 1, [10]dara.GeneralType{argInfo1}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_FCHDIR, syscallInfo)
	}
	if err := f.checkValid("chdir"); err != nil {
		return err
	}
	if e := f.pfd.Fchdir(); e != nil {
		return f.wrapErr("chdir", e)
	}
	return nil
}

// setDeadline sets the read and write deadline.
func (f *File) setDeadline(t time.Time) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[SetDeadline] : ")
		print(f.file.name)
		print(" ")
		print(t.String())
		argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type:dara.TIME, String: t.String()}
		retInfo := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_SETDEADLINE, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_SETDEADLINE, syscallInfo)
	}
	if err := f.checkValid("SetDeadline"); err != nil {
		return err
	}
	return f.pfd.SetDeadline(t)
}

// setReadDeadline sets the read deadline.
func (f *File) setReadDeadline(t time.Time) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[SetReadDeadline] : ")
		print(f.file.name)
		print(" ")
		print(t.String())
		argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type:dara.TIME, String: t.String()}
		retInfo := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_SETREADDEADLINE, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_SETREADDEADLINE, syscallInfo)
	}
	if err := f.checkValid("SetReadDeadline"); err != nil {
		return err
	}
	return f.pfd.SetReadDeadline(t)
}

// setWriteDeadline sets the write deadline.
func (f *File) setWriteDeadline(t time.Time) error {
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		print("[SetWriteDeadline] : ")
		print(f.file.name)
		print(" ")
		print(t.String())
		argInfo1 := dara.GeneralType{Type:dara.FILE, String: f.name}
		argInfo2 := dara.GeneralType{Type:dara.TIME, String: t.String()}
		retInfo := dara.GeneralType{Type:dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_SETWRITEDEADLINE, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_SETWRITEDEADLINE, syscallInfo)
	}
	if err := f.checkValid("SetWriteDeadline"); err != nil {
		return err
	}
	return f.pfd.SetWriteDeadline(t)
}

// checkValid checks whether f is valid for use.
// If not, it returns an appropriate error, perhaps incorporating the operation name op.
func (f *File) checkValid(op string) error {
	if f == nil {
		return ErrInvalid
	}
	return nil
}
