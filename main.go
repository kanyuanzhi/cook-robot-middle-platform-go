package main

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/boot"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/fxcron"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"log/slog"
	"time"
)

func main() {
	global.FXDb = boot.Sqlite()
	err := boot.InitSqliteDb()
	if err != nil {
		slog.Warn("数据表已存在")
		//return
	}
	//if global.FXDb == nil {
	//	err := boot.InitDb()
	//	if err != nil {
	//		slog.Error("初始化数据库失败")
	//		return
	//	}
	//}

	global.FXCommandRpcClient = boot.CommandRpcClient()
	global.FXSoftwareUpdaterRpcClient = boot.SoftwareUpdaterRpcClient()
	global.FXDataUpdaterRpcClient = boot.DataUpdaterRpcClient()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				go fxcron.FetchControllerStatus()
			}
		}
	}()

	boot.Cron()
	boot.Boot()
}
