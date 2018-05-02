package syscall

func Is_dara_profiling_on() bool {
    if v, ok := Getenv("DARA_PROFILING"); v == "" || !ok {
        return false
    }

    return true
}
