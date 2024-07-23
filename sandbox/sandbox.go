package sandbox

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"
)

type Result struct {
	Stdout   string
	Stderr   string
	CpuTime  int64
	RealTime int64
	Memory   int64
}

func run(cmd *exec.Cmd) (*Result, error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var usage syscall.Rusage
	start := time.Now()

	// Start process
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	pid := cmd.Process.Pid

	// Wait for process execution to end
	_, err := syscall.Wait4(pid, nil, 0, &usage)
	if err != nil {
		return nil, err
	}
	end := time.Now()

	result := &Result{
		CpuTime:  (usage.Utime.Nano() + usage.Stime.Nano()) / 1000000,
		RealTime: end.Sub(start).Milliseconds(),
		Memory:   usage.Maxrss,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}
	return result, nil
}
