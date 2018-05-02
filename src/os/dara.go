package os

func Is_dara_profiling_on() bool {
    if v := Getenv("DARA_PROFILING"); v != "" {
        return true
    }

    return false
}

