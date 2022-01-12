package process

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/sirupsen/logrus"
)

type Process struct {
	args                string
	cmd                 string
	env                 []string
	logOutput           bool
	outputLineDelimiter string
	pipeStderrToStdout  bool
	workingDir          string
	log                 *logrus.Logger
}

func (p *Process) logCommandOutput(command *cmd.Cmd) {
	ticker := time.NewTicker(50 * time.Millisecond)

	// Output our stderr and stdout
	go func() {
		if !p.logOutput {
			return
		}
		stdoutPosition := 0
		stderrPosition := 0
		for range ticker.C {
			if command.Status().Complete {
				ticker.Stop()
				return
			}
			fmt.Print(".")
			status := command.Status()
			nStdOut := len(status.Stdout)
			if nStdOut > stdoutPosition {
				line := status.Stdout[nStdOut-1]
				fmt.Println()
				p.log.Println(line)
				stdoutPosition = nStdOut
			}
			nStdErr := len(status.Stderr)
			if nStdErr > stderrPosition {
				line := status.Stderr[nStdErr-1]
				fmt.Println()
				if p.pipeStderrToStdout {
					p.log.Println(line)
				} else {
					p.log.Errorln(line)
				}
				stderrPosition = nStdErr
			}
		}

	}()
}

func (p *Process) run() (string, string, error) {
	command := cmd.NewCmd(p.cmd, strings.Split(p.args, " ")...)
	command.Dir = p.workingDir
	command.Env = append(os.Environ(), p.env...)
	p.log.Debugf("Running command %s %s", p.cmd, p.args)

	statusChan := command.Start()
	p.logCommandOutput(command)

	// Block waiting for command to exit
	finalStatus := <-statusChan
	p.log.Debugf("Exit code from process: %d, %v", finalStatus.Exit, finalStatus.Error)

	var err error
	if finalStatus.Exit != 0 {
		err = fmt.Errorf("%s failed with exit code %d", p.cmd, finalStatus.Exit)
	}

	return strings.Join(finalStatus.Stdout, "\n"),
		strings.Join(finalStatus.Stderr, "\n"),
		err
}

func New(cmd string, workingDir string, logger *logrus.Logger) *Process {
	return &Process{
		args:                "",
		cmd:                 cmd,
		logOutput:           true,
		log:                 logger,
		outputLineDelimiter: "",
		workingDir:          workingDir,
	}
}

func (p *Process) WithArgs(args string) *Process {
	p.args = args
	return p
}

func (p *Process) WithEnv(kvp string) *Process {
	p.env = append(p.env, kvp)
	return p
}

func (p *Process) WithSuppressedOutput() *Process {
	p.logOutput = false
	return p
}

func (p *Process) PipeStderrToStdout() *Process {
	p.pipeStderrToStdout = true
	return p
}

func (p *Process) Run() (string, string, error) {
	return p.run()
}
