package dara

type TypeNum int

const (
	INTEGER TypeNum = iota
	INTEGER64
	BOOL
	FLOAT
	STRING
	ARRAY
	ERROR
	POINTER
	FILE
	FILEINFO
	CONNECTION
	TIME
	PROCESS
	SIGNAL
	CONTEXT
	SOCKADDR
)

//Event Types
const (
	LOG_EVENT = iota
	SYSCALL_EVENT
	SEND_EVENT
	REC_EVENT
	SCHED_EVENT
)



//DaraProc is used to communicate control and data information between
//a single instrumented go runtime and the global scheduler. One of
//these structures is mapped into shared memory for each process that
//is launched during an execution. If there are more runtimes active
//Than DaraProcs the additional runtimes will not be controlled by the
//global scheduler, or Segfault immediately 
type DaraProc struct {

	//Lock is used to control the execution of a process. A process
	//which is running but Not scheduled will spin on this lock using
	//checkandset operations If the lock is held The owner can modify
	//the state of the DaraProc
	Lock uint32

	//SyscallLock is used to control the reporting of the syscalls.
	SyscallLock uint32

	//Run is a deprecated var with multiple purposes. Procs set their
	//Run to -1 when they Are done running (in replay mode) to let the
	//scheduler know they are done. The global scheduler sets this
	//value to 2 to let the runtime know its replay, and 3 for record
	//1 is used to denote the first event, and 0 indicates this
	//variable has not been initialized Originally Run was intended to
	//report the id of the goroutine that was executed, but that was
	//not always the same so the program counter  was needed, now
	//RunningRoutine is used to report this. The global scheduler sets
	// this to -4 to inform the local schedulers that replay is ended
	Run int
	// Syscall number at which the running routine is blocked on. -1 means that
	// there is no syscall on which the daraproc is blocked
	Syscall int
	//RunningRoutine is the goroutine scheduled, running, or ran, for
	//any single replayed event in a schedule. In Record, the
	//executed goroutine is reported via this variable, in Replay the
	//global scheduler tells the runtime which routine to run with
	//RunningRoutine
	RunningRoutine RoutineInfo
	//Routines is the full set of goroutines a process is allowed to
	//run. The total number is allocated upfront so that shared memory
	//does not need to be resized dynamically. After each iteration
	//of scheduling runtimes update the states of all their routines
	//via this structure
	Routines [MAXGOROUTINES]RoutineInfo
	//TODO document
	Epoch int
	LogIndex int
	Log [MAXLOGENTRIES]EncEvent
}

//RoutineInfo contains data specific to a single goroutine
type RoutineInfo struct {
        //Set to one of the statuses in the constant block above
        Status uint32
        //Goroutine id as set by the runtime. This is sometimes usefull
        //for detecting which routine is which, but it is not always the
        //same between runs. However 1 is allways main, while 2 is a
        //finalizer, and 3 is a garbage collection invocator.
        Gid int
        //Program counter that this goroutine was launched from
        Gpc uintptr
        //A count of how many other goroutines were launched on the same
        //pc prior to this goroutine. (Gpc,Routinecount) is a unique id
        //for a goroutine on a given processor.
        RoutineCount int
        //A textual description of the function this goroutine was forked
        //from.In the future it can be removed.
        FuncInfo [64]byte
}

type EncEvent struct {
	Type int
	P int
	G RoutineInfo
	Epoch int
	ELE EncLogEntry
	SyscallInfo GeneralSyscall
	EM EncodedMessage
}

type EncLogEntry struct {
	Length int
	LogID [VARBUFLEN]byte
	Vars [MAXLOGVARIABLES] EncNameValuePair
}

type GeneralType struct {
	Type TypeNum
	Integer int
	Bool bool
	Float float32
	Integer64 int64
	String string
	Unsupported rune
}

type GeneralSyscall struct {
	SyscallNum int
	NumArgs int
	NumRets int
	Args [10]GeneralType
	Rets [10]GeneralType
}

//TODO fill out this structure
type EncodedMessage struct {
	Body [VARBUFLEN]byte
}

type EncNameValuePair struct {
	VarName [VARBUFLEN]byte
	Value [VARBUFLEN]byte
	Type [VARBUFLEN]byte
}

type Event struct {
	Type int
	P int
	G RoutineInfo
	Epoch int
	LE LogEntry
	SyscallInfo GeneralSyscall
	Msg Message
}

type LogEntry struct {
	LogID string
	Vars []NameValuePair
}

type NameValuePair struct {
	VarName string
	Value interface{}
	Type string
}

type Message struct {
	Body string
}


type DaraProcStatus uint32

// The numbering has to match the Goroutine states from runtime/proc.go
const (
        Idle DaraProcStatus = iota // 0
        Runnable // 1
        Running // 2
        Syscall // 3
        Waiting // 4
        Moribound_Unused // 5
        Dead // 6
        Enqueue_Unused // 7
        Copystack // 8
        Scan DaraProcStatus = 0x1000
        ScanRunnable = Scan + Runnable // 0x1001
        ScanRunning  = Scan + Running  // 0x1002
        ScanSyscall  = Scan + Syscall  // 0x1003
        ScanWaiting  = Scan + Waiting  // 0x1004
)

func GetDaraProcStatus(status uint32) DaraProcStatus {
	return DaraProcStatus(status)
}

//Type which encapsulates a single schedule
//TODO integerate with vaastav to build a single schedule for DPOR
type Schedule []Event

