package boot

import (
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/config"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/model"
	"gorm.io/gorm"
	"log/slog"
)

func Sqlite() *gorm.DB {
	fxSqliteConfig := global.FXConfig.Sqlite
	sqliteDialector := config.SqliteDialector(fxSqliteConfig)
	gormConfig := config.GormConfig()
	if db, err := gorm.Open(sqliteDialector, &gormConfig); err != nil {
		slog.Error("Connect sqlite failed")
		return nil
	} else {
		slog.Info("Connect sqlite success")
		return db
	}
}

var migrateList = []interface{}{
	model.SysDish{},
	model.SysCuisine{},
	model.SysSeasoning{},
	model.SysIngredient{},
	model.SysIngredientType{},
	model.SysIngredientShape{},
}

func InitSqliteDb() error {
	global.FXDb = Sqlite()
	err := global.FXDb.AutoMigrate(migrateList...)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func IsSqliteInit() bool {
	return global.FXDb != nil
}
