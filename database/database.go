package database

import (
	"os"
	"path/filepath"

	"github.com/theamniel/scheduler/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(OpenConnection(), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&types.Schedule{})
	return db, nil
}

func OpenConnection() gorm.Dialector {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir, _ := filepath.Split(executable)
	return sqlite.Open(dir + "scheduler.db")
}
