package exit

import "os"

const (
	ExitSuccess      = 0
	ExitRepoError    = 1
	ExitHTTPSrvError = 2
	ExitOtherError   = 3
)

// Code тип для передчи информации о типе выхода из приложения
type Code int

// Exit метод для организации завершения приложения
func (c Code) Exit() {
	os.Exit(int(c))
}
