package exit

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestCode_Exit(t *testing.T) {
	if os.Getenv("TEST_EXIT_CRASHER") == "1" {
		codeStr := os.Getenv("TEST_EXIT_CODE")
		codeInt, _ := strconv.Atoi(codeStr)

		Code(codeInt).Exit()
		return
	}

	tests := []struct {
		name         string
		inputCode    int
		expectedCode int
	}{
		{"Success", ExitSuccess, 0},
		{"RepoError", ExitRepoError, 1},
		{"HTTPError", ExitHTTPSrvError, 2},
		{"OtherError", ExitOtherError, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmd := exec.Command(os.Args[0], "-test.run=TestCode_Exit")

			cmd.Env = append(os.Environ(),
				"TEST_EXIT_CRASHER=1",
				"TEST_EXIT_CODE="+strconv.Itoa(tt.inputCode),
			)

			err := cmd.Run()

			if tt.expectedCode == 0 {
				if err != nil {
					t.Errorf("expected exit code 0, got error: %v", err)
				}
				return
			}

			if e, ok := err.(*exec.ExitError); ok {
				if e.ExitCode() != tt.expectedCode {
					t.Errorf("expected exit code %d, got %d", tt.expectedCode, e.ExitCode())
				}
			} else {
				t.Errorf("expected exit error, got %v", err)
			}
		})
	}
}
