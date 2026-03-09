package main

import (
	"fmt"
	logger "github.com/DerrickHsu97/xlog4go"
)

func main() {
	// 使用多 logger 配置文件初始化
	if err := logger.SetupMultiLogWithConf("./multi-log.json"); err != nil {
		panic(err)
	}
	defer logger.Close()

	fmt.Println("========== 默认 Logger (业务日志) ==========")
	// GetLogger() 获取默认 logger
	logger.Info("这是一条业务 info 日志")
	logger.Debug("这条 debug 日志不会显示，因为默认 logger 级别是 info")
	logger.Error("这是一条业务 error 日志")

	fmt.Println("\n========== GORM Logger ==========")
	// GetLoggerByName("gorm") 获取 gorm 专用 logger
	// 如果 gorm logger 不存在，会返回默认 logger
	gormLog := logger.GetLoggerByName("gorm")
	gormLog.Debug("SELECT * FROM users WHERE id = 1")
	gormLog.Info("gorm 执行了 SQL 查询")

	fmt.Println("\n========== Redis Logger ==========")
	redisLog := logger.GetLoggerByName("redis")
	redisLog.Debug("GET key: user:1001")
	redisLog.Info("redis 连接成功")

	fmt.Println("\n========== Access Logger ==========")
	accessLog := logger.GetLoggerByName("access")
	accessLog.Info("POST /api/login 200 15ms")
	accessLog.Info("GET /api/user 200 8ms")

	fmt.Println("\n========== 不存在的 Logger ==========")
	// 获取不存在的 logger，会返回默认 logger
	unknownLog := logger.GetLoggerByName("unknown")
	unknownLog.Info("这条日志会输出到默认 logger，因为 'unknown' 不存在")

	fmt.Println("\n========== 手动创建和注册 Logger ==========")
	// 也可以手动创建 logger 并注册
	manualLog := logger.NewLoggerWithName("manual")
	manualLog.SetLevel(logger.DEBUG)
	fw := logger.NewFileWriter()
	fw.SetFileName("./logs/manual.log")
	fw.SetPathPattern("./logs/manual.log.%Y%m%d")
	fw.SetLogLevelFloor(logger.TRACE)
	fw.SetLogLevelCeil(logger.ERROR)
	manualLog.Register(fw)

	retrievedLog := logger.GetLoggerByName("manual")
	retrievedLog.Info("这是手动创建的 logger 输出的日志")
}
