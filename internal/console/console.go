package console

type Console interface {
	Log(level LogLevel, text string)
	Print(message string)
	PrintFlockTitle()
	PrintBox(nodeInfo NodeInfo)
}

type NodeInfo struct {
	NodeID          string
	NodeName        string
	IsBootstrapNode bool
	GossipPort      int
	ApiPort         int
	Protocol        string
	Status          string
}

type LogLevel uint8

const (
	LogLevelDebug LogLevel = iota + 1
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func New() Console {
	return PtermConsole{}
}
