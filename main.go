package main

import (
	"github.com/kanyuanzhi/middle-platform/boot"
	"github.com/kanyuanzhi/middle-platform/global"
	"log/slog"
)

func main() {
	global.FXDb = boot.Sqlite()
	err := boot.InitSqliteDb()
	if err != nil {
		slog.Warn("数据表已存在")
		return
	}
	//if global.FXDb == nil {
	//	err := boot.InitDb()
	//	if err != nil {
	//		slog.Error("初始化数据库失败")
	//		return
	//	}
	//}

	boot.Boot()
}
