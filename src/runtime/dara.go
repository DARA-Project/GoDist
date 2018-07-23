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

func Report_Syscall_To_Scheduler(syscallID int, syscallInfo dara.GeneralSyscall) {
	report_syscall(syscallID, syscallInfo)
}
