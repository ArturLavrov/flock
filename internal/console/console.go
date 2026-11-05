package console

type Console interface {
	Log(level LogLevel, text string)
	Print(message string)
	PrintFlockTitle()
	PrintAgentInfo(agentID, serviceName string, gossipPort, apiPort int, protocol, version, status string)
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
