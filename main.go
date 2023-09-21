package main

import (
	"github.com/kanyuanzhi/middle-platform/boot"
	"github.com/kanyuanzhi/middle-platform/fxcron"
	"github.com/kanyuanzhi/middle-platform/global"
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
