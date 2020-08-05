package runtime

import (
	"dara"
	"unsafe"
)

func Is_dara_profiling_on() bool {
	if v := gogetenv("DARA_PROFILING"); v == "" {
		return false
	}
	return true
}

func Is_Dara_On() bool {
	if v := gogetenv("DARAON"); v == "" {
		return false
	}

	return true
}

func ReportBlockCoverage(blockID string) {
    // Update counter for the block
    if DaraInitialised {
		var start int64
		if Microbenchmark || Nanobenchmark {
			start = nanotime()
		}
        if v, ok := CoverageInfo[blockID]; ok {
            CoverageInfo[blockID] = v + 1
        } else {
            CoverageInfo[blockID] = uint64(1)
		}
		if Microbenchmark || Nanobenchmark {
			end := nanotime() - start
			println("COVERAGE,", end)
		}
        dprint(dara.DEBUG, func() {println("Reporting coverage for block:", blockID)} )
    }
}

func Report_Syscall_To_Scheduler(syscallID int, syscallInfo dara.GeneralSyscall) {
	//report_syscall(syscallID, syscallInfo) //TODO remove this it's redundent and slow
	var start int64
	if Microbenchmark || Nanobenchmark {
		// Get the time in nanoseconds (atleast I think this is what the function does)
		start = nanotime()
	}
	LogSyscall(syscallInfo)
	if Microbenchmark || Nanobenchmark {
		end := nanotime() - start
		println("SYSCALL,", end)
	}
}

func Dara_Debug_Print(pfunc func()) {
    dprint(dara.DEBUG,pfunc)
}

func dara_Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Functions for reporting data to users

type emptyInterface struct {
	typ *_type
	word unsafe.Pointer
}

func NumDeliveries(ch interface{}) int {
	if DaraInitialised {
		ef := (*emptyInterface)(unsafe.Pointer(&ch))
		return ChanRecvInfo[ef.word]
	}
	return -1
}

func NumSendings(ch interface{}) int {
	if DaraInitialised {
		ef := (*emptyInterface)(unsafe.Pointer(&ch))
		return ChanSendInfo[ef.word]
	}
	return -1
}

func DaraProcessID() int {
	if DaraInitialised {
		return DPid
	}
	return -1
}