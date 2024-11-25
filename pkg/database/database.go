package database

import (
	"Yakudza/pkg/config"
	"Yakudza/pkg/logger"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"time"
)

var dbase *gorm.DB
var g_cfg *config.Config

// Init - Инициализация базы данных
func Init(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Database.Host,
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log.Default.LogMode(log.Silent),
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	g_cfg = cfg

	return db, nil
}

// GetDB - Получение ссылки на экземпляр базы данных
func GetDB() *gorm.DB {
	if dbase == nil {
		dbase, _ = Init(g_cfg)
		sleep := time.Duration(1)
		for dbase == nil {
			sleep += 1
			logger.Info("Не удалось подключиться к базе данных, повторное подключение через %d секунд", sleep)
			time.Sleep(sleep * time.Second)
			dbase, _ = Init(g_cfg)
			if sleep >= 30 {
				sleep = time.Duration(1)
			}
		}
	}
	return dbase
}

func CloseConnection() {
	db := GetDB()
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	fmt.Println("Close database connections")
}
