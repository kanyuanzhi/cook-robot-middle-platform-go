package boot

import (
	"database/sql"
	"fmt"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/config"
	"github.com/kanyuanzhi/cook-robot-middle-platform-go/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

func Postgres() *gorm.DB {
	fxPostgresConfig := global.FXConfig.Postgres
	postgresConfig := config.PostgresConfig(fxPostgresConfig)
	gormConfig := config.GormConfig()
	if db, err := gorm.Open(postgres.New(postgresConfig), &gormConfig); err != nil {
		slog.Error("Connect postgres failed")
		return nil
	} else {
		slog.Info("Connect postgres success")
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("Get postgres instance failed")
			return nil
		}
		sqlDB.SetMaxOpenConns(10) // 最大打开连接数
		sqlDB.SetMaxIdleConns(5)  // 最大闲置连接数
		return db
	}
}

func InitDb() error {
	if err := createDatabase(); err != nil {
		return err
	}
	global.FXDb = Postgres()
	err := global.FXDb.AutoMigrate(migrateList...)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func createDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		global.FXConfig.Postgres.Host, global.FXConfig.Postgres.User, global.FXConfig.Postgres.Password, global.FXConfig.Postgres.Port)
	createSql := fmt.Sprintf("CREATE DATABASE %s WITH ENCODING UTF8", global.FXConfig.Postgres.Database)
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(sqlDB)
	if err = sqlDB.Ping(); err != nil {
		slog.Error(err.Error())
		return err
	}
	if _, err = sqlDB.Exec(createSql); err != nil {
		slog.Error(err.Error())
	}
	return err
}
