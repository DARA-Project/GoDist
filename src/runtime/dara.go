package runtime

import (
	"dara"
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
        if v, ok := CoverageInfo[blockID]; ok {
            CoverageInfo[blockID] = v + 1
        } else {
            CoverageInfo[blockID] = uint64(1)
        }
        dprint(dara.DEBUG, func() {println("Reporting coverage for block:", blockID)} )
    }
}

func Report_Syscall_To_Scheduler(syscallID int, syscallInfo dara.GeneralSyscall) {
	//report_syscall(syscallID, syscallInfo) //TODO remove this it's redundent and slow
	LogSyscall(syscallInfo)
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