package xlog4go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigConsole(t *testing.T) {
	if err := SetupLogWithConf("./test/test-console.json"); err != nil {
		assert.NoError(t, err)
	}
	Trace("trace")
	Debug("debug")
	Info("info")
	Warn("warning")
	Error("error")
	Fatal("fatal")
}

func TestConfigFile(t *testing.T) {
	if err := SetupLogWithConf("./test/test.json"); err != nil {
		assert.NoError(t, err)
	}
}

func TestConfigErrFile(t *testing.T) {
	if err := SetupLogWithConf("./test/test-err.json"); err != nil {
		assert.Error(t, err)
	}
}

func TestConfigNotExistFile(t *testing.T) {
	if err := SetupLogWithConf("./test/test-exist.json"); err != nil {
		assert.Error(t, err)
	}
}
