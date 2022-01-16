package process

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func createProcess(cmd string, workingDir string) *Process {
	logger, _ := test.NewNullLogger()
	return New(cmd, workingDir, logger)
}

func TestNew(t *testing.T) {
	p := createProcess("echo", "./wd")

	assert.Equal(t, "", p.args)
	assert.Equal(t, "echo", p.cmd)
	assert.Equal(t, true, p.logOutput)
	assert.Equal(t, "", p.outputLineDelimiter)
	assert.Equal(t, "./wd", p.workingDir)
}

func TestRunWithArgs(t *testing.T) {
	stdOut, stdErr, err := createProcess("echo", "").
		WithArgs("abc").
		Run()

	assert.Nil(t, err)
	assert.Equal(t, "", stdErr)
	assert.Equal(t, "abc", stdOut)
}

func TestRunWithSuppressedOutput(t *testing.T) {
	logger, logHook := test.NewNullLogger()
	logger.Level = logrus.InfoLevel
	New("echo", "", logger).
		WithArgs("abc").
		WithSuppressedOutput().
		Run()

	logEntries := logHook.AllEntries()
	assert.Equal(t, 0, len(logEntries), "output was logged when it was not expected")
}
