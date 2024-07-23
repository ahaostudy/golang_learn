package sandbox

import (
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	BaseURL = "exe"
)

func testRun(t *testing.T, cmd *exec.Cmd) {
	result, err := run(cmd)
	if err != nil {
		panic(err)
	}
	t.Logf("result: %#v\n", result)
}

func TestRunGolang(t *testing.T) {
	testRun(t, exec.Command(filepath.Join(BaseURL, "golang/main")))
}

func TestRunCPP(t *testing.T) {
	testRun(t, exec.Command(filepath.Join(BaseURL, "cpp/main")))
}
