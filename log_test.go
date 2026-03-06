package dlog4go

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalLogger(t *testing.T) {
	public := NewFileWriter()
	public.SetLogLevelFloor(ERROR)
	public.SetLogLevelCeil(PUBLIC)

	buf := &bytes.Buffer{}
	public.fileBufWriter = bufio.NewWriter(buf)

	logger_default.writers = []Writer{public}
	logger_default.SetLevel(WARNING)
	logger_default.SetLayout("2010")

	Trace("trace")
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	Fatal("fatal")
	Public("public")
	Close()

	assert.NotContains(t, buf.String(), "warn")
	assert.NotContains(t, buf.String(), "debug")
	assert.NotContains(t, buf.String(), "trace")

	assert.Contains(t, buf.String(), "error")
	assert.Contains(t, buf.String(), "fatal")
	assert.Contains(t, buf.String(), "public")
}

func TestLogger(t *testing.T) {
	public := NewFileWriter()
	public.SetLogLevelFloor(ERROR)
	public.SetLogLevelCeil(PUBLIC)

	buf := &bytes.Buffer{}
	public.fileBufWriter = bufio.NewWriter(buf)

	defaut := takeup
	takeup = true
	log := NewLogger()
	takeup = defaut

	log.writers = []Writer{public}
	log.SetLevel(WARNING)
	log.SetLayout("2010")

	log.Trace("trace")
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	log.Fatal("fatal")
	log.Public("public")
	log.Close()

	assert.NotContains(t, buf.String(), "warn")
	assert.NotContains(t, buf.String(), "debug")
	assert.NotContains(t, buf.String(), "trace")

	assert.Contains(t, buf.String(), "error")
	assert.Contains(t, buf.String(), "fatal")
	assert.Contains(t, buf.String(), "public")
}
