package main

import (
	logger "git.xiaojukeji.com/shield-arch/dlog4go"

	"time"
)

func main() {
	if err := logger.SetupLogWithConf("./log.json"); err != nil {
		panic(err)
	}
	defer logger.Close()

	var name = "hellow world"
	for {
		logger.Trace("dlog4go by %s", name)
		logger.Debug("dlog4go by %s", name)
		logger.Info("dlog4go by %s", name)
		logger.Warn("dlog4go by %s", name)
		logger.Error("dlog4go by %s", name)
		logger.Fatal("dlog4go by %s", name)
		logger.Public("dlog4go by %s", name)

		time.Sleep(time.Second * 1)
	}
}
