package store

import (
	"sync"

	"github.com/lfcypo/scrkun/logger"
	"github.com/lfcypo/viperx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var databaseLog = logger.New("Database")

var db *gorm.DB
var doOnceInitDb = &sync.Once{}

func initDatabase() {
	doOnceInitDb.Do(func() {
		_initDatabase()
	})
}

func _initDatabase() {
	dsn := viperx.GetString("db.path", "data.db")
	dial := sqlite.Open(dsn)
	config := &gorm.Config{}
	dbInner, err := gorm.Open(dial, config)
	if err != nil {
		databaseLog.Errorf("初始化数据库失败: %v", err)
	}

	db = dbInner

	err = autoMigrate()
	if err != nil {
		databaseLog.Errorf("自动迁移数据库失败: %v", err)
	}
}

func GetDatabase() *gorm.DB {
	initDatabase()
	return db
}
