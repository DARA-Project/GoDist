// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

// Unix environment variables.

package syscall

import "dara"
import "runtime"
import "sync"

var (
	// envOnce guards initialization by copyenv, which populates env.
	envOnce sync.Once

	// envLock guards env and envs.
	envLock sync.RWMutex

	// env maps from an environment variable to its first occurrence in envs.
	env map[string]int

	// envs is provided by the runtime. elements are expected to
	// be of the form "key=value". An empty string means deleted
	// (or a duplicate to be ignored).
	envs []string = runtime_envs()
)

func runtime_envs() []string // in package runtime

// setenv_c and unsetenv_c are provided by the runtime but are no-ops
// if cgo isn't loaded.
func setenv_c(k, v string)
func unsetenv_c(k string)

func copyenv() {
	env = make(map[string]int)
	for i, s := range envs {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				key := s[:j]
				if _, ok := env[key]; !ok {
					env[key] = i // first mention of key
				} else {
					// Clear duplicate keys. This permits Unsetenv to
					// safely delete only the first item without
					// worrying about unshadowing a later one,
					// which might be a security problem.
					envs[i] = ""
				}
				break
			}
		}
	}
}

func Unsetenv(key string) error {
	envOnce.Do(copyenv)

	envLock.Lock()
	defer envLock.Unlock()

	if i, ok := env[key]; ok {
		envs[i] = ""
		delete(env, key)
	}
	unsetenv_c(key)
	// DARA Instrumentation
	if (runtime.Is_dara_profiling_on()) {
		println("[UNSETENV] : " + key)
		argInfo := dara.GeneralType{Type: dara.STRING, String:key}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported : dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_UNSETENV, 1, 1, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_UNSETENV, syscallInfo)
	}
	return nil
}

func Getenv(key string) (value string, found bool) {
	argInfo := dara.GeneralType{Type: dara.STRING, String:key}
	envOnce.Do(copyenv)
	if len(key) == 0 {
		return "", false
	}

	envLock.RLock()
	defer envLock.RUnlock()

	i, ok := env[key]
	if !ok {
		// DARA Instrumentation
		if (runtime.Is_dara_profiling_on()) {
			println("[GETENV] : " + key)
			retInfo1 := dara.GeneralType{Type: dara.STRING, String:""}
			retInfo2 := dara.GeneralType{Type: dara.BOOL, Bool:false}
			syscallInfo := dara.GeneralSyscall{dara.DSYS_GETENV, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
			runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETENV, syscallInfo)
		}
		return "", false
	}
	s := envs[i]
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			// DARA Instrumentation
			if (runtime.Is_dara_profiling_on()) {
				println("[GETENV] : " + key)
				retInfo1 := dara.GeneralType{Type: dara.STRING, String:""}
				retInfo2 := dara.GeneralType{Type: dara.BOOL, Bool:false}
				syscallInfo := dara.GeneralSyscall{dara.DSYS_GETENV, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
				runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETENV, syscallInfo)
			}
			return s[i+1:], true
		}
	}
	// DARA Instrumentation
	if (runtime.Is_dara_profiling_on()) {
		println("[GETENV] : " + key)
		retInfo1 := dara.GeneralType{Type: dara.STRING, String:""}
		retInfo2 := dara.GeneralType{Type: dara.BOOL, Bool:false}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETENV, 1, 2, [10]dara.GeneralType{argInfo}, [10]dara.GeneralType{retInfo1, retInfo2}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_GETENV, syscallInfo)
	}
	return "", false
}

func Setenv(key, value string) error {
	// DARA Instrumentation
	if (runtime.Is_dara_profiling_on()) {
		println("[SETENV] : " + key +  " "  + value)
		argInfo1 := dara.GeneralType{Type: dara.STRING, String:key}
		argInfo2 := dara.GeneralType{Type: dara.STRING, String:value}
		retInfo := dara.GeneralType{Type: dara.ERROR, Unsupported : dara.UNSUPPORTEDVAL}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_GETENV, 2, 1, [10]dara.GeneralType{argInfo1, argInfo2}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_SETENV, syscallInfo)
	}
	envOnce.Do(copyenv)
	if len(key) == 0 {
		return EINVAL
	}
	for i := 0; i < len(key); i++ {
		if key[i] == '=' || key[i] == 0 {
			return EINVAL
		}
	}
	for i := 0; i < len(value); i++ {
		if value[i] == 0 {
			return EINVAL
		}
	}

	envLock.Lock()
	defer envLock.Unlock()

	i, ok := env[key]
	kv := key + "=" + value
	if ok {
		envs[i] = kv
	} else {
		i = len(envs)
		envs = append(envs, kv)
	}
	env[key] = i
	setenv_c(key, value)
	return nil
}

func Clearenv() {
	envOnce.Do(copyenv) // prevent copyenv in Getenv/Setenv

	envLock.Lock()
	defer envLock.Unlock()

	for k := range env {
		unsetenv_c(k)
	}
	env = make(map[string]int)
	envs = []string{}
	// DARA Instrumentation
	if (runtime.Is_dara_profiling_on()) {
		println("[CLEARENV]")
		syscallInfo := dara.GeneralSyscall{dara.DSYS_CLEARENV, 0, 0, [10]dara.GeneralType{}, [10]dara.GeneralType{}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_CLEARENV, syscallInfo)
	}
}

func Environ() []string {
	envOnce.Do(copyenv)
	envLock.RLock()
	defer envLock.RUnlock()
	a := make([]string, 0, len(envs))
	for _, env := range envs {
		if env != "" {
			a = append(a, env)
		}
	}
	// DARA Instrumentation
	if (runtime.Is_dara_profiling_on()) {
		println("[ENVIRON]")
		retInfo := dara.GeneralType{Type: dara.ARRAY, Integer: len(a)}
		syscallInfo := dara.GeneralSyscall{dara.DSYS_ENVIRON, 0, 1, [10]dara.GeneralType{}, [10]dara.GeneralType{retInfo}}
		runtime.Report_Syscall_To_Scheduler(dara.DSYS_ENVIRON, syscallInfo)
	}
	return a
}
