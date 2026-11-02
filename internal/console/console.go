package console

type Console interface {
	Print(message string)
}

func New() Console {
	return PtermConsole{}
}
