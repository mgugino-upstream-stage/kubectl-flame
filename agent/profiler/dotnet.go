package profiler

import (
	"bytes"
	"github.com/VerizonMedia/kubectl-flame/agent/details"
	"github.com/VerizonMedia/kubectl-flame/agent/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	containertmp = "/app/containertmp.sh"
	dntrace      = "/app/dotnet-trace"
	flameout     = "/tmp/trace"
	ssfile       = "/tmp/trace.speedscope.json"
)

type DotNetProfiler struct{}

func (p *DotNetProfiler) SetUp(job *details.ProfilingJob) error {
	return nil
}

func (p *DotNetProfiler) Invoke(job *details.ProfilingJob) error {

	// Determine real path to container tmp filesystem.
	cid := strings.Split(job.ContainerID, "://")[1]
	cmd := exec.Command(containertmp, cid)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	// TODO: better error handling here if path parsing fails.
	realTMP := out.String()

	// Create symlink
	err = os.Symlink(realTMP, "/app/tmp")
	if err != nil {
					return err
	}
	// Run dotnet trace
	// /app/dotnet-trace collect -p 1
	cmd2 := exec.Command(dntrace, "collect", "--format", "Speedscope", "-p", "1", "-o", flameout)
	cmd2.Env = os.Environ()
	cmd2.Env = append(cmd.Env, "TMPDIR=/app/tmp")

	err = cmd2.Start()
	if err != nil {
		return err
	}

	// sleep duration seconds
	time.Sleep(job.Duration)
	// Send sigint
	err = cmd2.Process.Signal(os.Interrupt)
	// TODO: catch return signal of dntrace
	if err != nil {
		// If the above fails, we'll leak the process, but we're stopping anyway.
		return err
	}

	// TODO: maybe set a timeout channel for this.
	cmd2.Wait()

	return utils.PublishFlameGraph(ssfile)
}
