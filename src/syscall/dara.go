package syscall

import (
	"runtime"
	"strconv"
	"unsafe"
)

// Duplication across runtime, godist scheduler and the go source code

type DaraProc struct {
	Lock int32
	Run int
	TotalRoutines int
	RunningRoutine RoutineInfo
	Routines [MAXGOROUTINES]RoutineInfo
}

type RoutineInfo struct {
	Status uint32
	Gid int
	Gpc uintptr
	RoutineCount int
	FuncInfo [64]byte
	Syscall int
}

//DARA Specific consts
const (
	//TODO TODO shared constants from the global scheduler, document
	//all of this well
	CHANNELS = 5 //TODO should be the same as procs
	DARAFD = 666
	UNLOCKED = 0
	LOCKED = 1
	SCHEDLEN = 1000000000
	PROCS = 3
	MAXGOROUTINES = 4096

	PAGESIZE = 4096
	SHAREDMEMPAGES = 65536

	//debug levels
	DEBUG = iota
	INFO
	WARN
	FATAL
	OFF
)

//These Constants are the arguments for MMAP, not all of the arguments
//are used, but are here for completeness
const(
        _EINTR  = 0x4
        _EAGAIN = 0xb
        _ENOMEM = 0xc

        _PROT_NONE  = 0x0
        _PROT_READ  = 0x1
        _PROT_WRITE = 0x2
        _PROT_EXEC  = 0x4

        _MAP_ANON    = 0x20
        _MAP_PRIVATE = 0x2
        _MAP_FIXED   = 0x10

        _MAP_SHARED = 0x01
)


// End duplication

var (
	isDaraInitialized bool
	p unsafe.Pointer //unsafe shared memory pointer
	err int
	procchan *[CHANNELS]DaraProc
	DPid int
)

func Is_dara_profiling_on() bool {
    if v, ok := Getenv("DARA_PROFILING"); v == "" || !ok {
        return false
    }

    return true
}

func Is_Dara_On() bool {
	if v, ok := Getenv("DARAON"); v == "" || !ok {
		return false
	}

	return true
}

func get_Dara_Pid() int {
	if val, ok := Getenv("DARAPID"); ok {
		if pid, err := strconv.Atoi(val); err == nil {
			return int(pid)
		}
	}
	// TODO Make this fatal
	println("DARA turned on but DARAPID not set")
	return -1
}

func initDara() {
	p, err = runtime.Mmap(nil, SHAREDMEMPAGES*PAGESIZE,_PROT_READ|_PROT_WRITE, _MAP_SHARED,DARAFD,0)

	if err != 0 {
		// TODO : Make this fatal
		println(err)
	}

	procchan = (*[CHANNELS]DaraProc)(p)
	isDaraInitialized = true
}

func Report_Syscall_To_Scheduler(syscallID int) {
	println("Inside syscall reporter to report syscall #",syscallID)
	if Is_Dara_On() {
		if !isDaraInitialized {
			initDara()
		}
	}
	DPid = get_Dara_Pid()
	if DPid != -1 {
		println("Reporting syscall #",syscallID)
		procchan[DPid].RunningRoutine.Syscall = syscallID
	}
}

