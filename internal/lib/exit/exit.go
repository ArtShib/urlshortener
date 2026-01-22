package exit

import "os"

const (
	ExitSuccess      = 0
	ExitRepoError    = 1
	ExitHTTPSrvError = 2
	ExitOtherError   = 3
)

type Code int

func (c Code) Exit() {
	os.Exit(int(c))
}
