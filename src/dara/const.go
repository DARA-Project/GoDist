package dara

//DARA Specific consts
//Constants for shared memory.
const (
	MAP_SHARED = 0x01 // used for mmap (not defined in linux defs)
	PROT_READ = 0x1
	PROT_WRITE = 0x2

	//The number of preallocated communication channels in shared
	//memory. This value can change in the future, for now 5 is the
	//maximum number of Processes. Invariant: dara.CHANNELS >= procs. TODO
	//assert this
	CHANNELS = 4

	//File discriptor for shared memory. This is set in the runscript.
	DARAFD = 666

	//State of spin locks. These are used by cas operations to control
	//the execution of the insturmented runtimes
	UNLOCKED = 0
	LOCKED = 1

	//The total size of the shared memory region is
        //dara.PAGESIZE*dara.SHAREDMEMPAGES
	PAGESIZE = 4096
	SHAREDMEMPAGES = 65536*4

	// TODO : This must be automated and not hardcoded
	DARAPROCSIZE = 106791048

	SCHEDLEN = 1000000000
	PROCS = 3
	MAXGOROUTINES = 4096

	MAXLOGENTRIES = 4096
	MAXLOGVARIABLES = 128
	VARBUFLEN = 64
	UNSUPPORTEDVAL = 2440

	//debug levels
	DEBUG = iota
	INFO
	WARN
	FATAL
	OFF
)

//loggging consts
const (
	BOOL_STRING = "bool"
	INT_STRING = "int"
	FLOAT_STRING = "float"
	STRING_STRING = "string"
)

//Goroutine states from runtime/proc.go
const (
	_Gidle = iota // 0
	_Grunnable // 1
	_Grunning // 2
	_Gsyscall // 3
	_Gwaiting // 4
	_Gmoribund_unused // 5
	_Gdead // 6
	_Genqueue_unused // 7
	_Gcopystack // 8
	_Gscan         = 0x1000
	_Gscanrunnable = _Gscan + _Grunnable // 0x1001
	_Gscanrunning  = _Gscan + _Grunning  // 0x1002
	_Gscansyscall  = _Gscan + _Gsyscall  // 0x1003
	_Gscanwaiting  = _Gscan + _Gwaiting  // 0x1004
)

//array of status to string from runtime/proc.go
var GStatusStrings = [...]string{
        _Gidle:      "idle",
        _Grunnable:  "runnable",
        _Grunning:   "running",
        _Gsyscall:   "syscall",
        _Gwaiting:   "waiting",
        _Gdead:      "dead",
        _Gcopystack: "copystack",
}

