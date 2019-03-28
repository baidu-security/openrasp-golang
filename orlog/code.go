package orlog

type ModuleCode uint32

const (
	Log       ModuleCode = 20002
	Config    ModuleCode = 20004
	Plugin    ModuleCode = 20005
	Runtime   ModuleCode = 20006
	Register  ModuleCode = 20008
	Heartbeat ModuleCode = 20009
)
