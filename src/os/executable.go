// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"dara"
	"runtime"
)

// Executable returns the path name for the executable that started
// the current process. There is no guarantee that the path is still
// pointing to the correct executable. If a symlink was used to start
// the process, depending on the operating system, the result might
// be the symlink or the path it pointed to. If a stable result is
// needed, path/filepath.EvalSymlinks might help.
//
// Executable returns an absolute path unless an error occurred.
//
// The main use case is finding resources located relative to an
// executable.
//
// Executable is not supported on nacl.
func Executable() (string, error) {
	str, err := executable()
	// DARA Instrumentation
	if runtime.Is_dara_profiling_on() {
		println("[EXECUTABLE]")
		retInfo1 := dara.GeneralType{Type: dara.STRING, String: str}
		retInfo2 := dara.GeneralType{Type: dara.ERROR, Unsupported: dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_EXECUTABLE, 0, 2, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_EXECUTABLE, syscallInfo)
	}
	return str, err
}
