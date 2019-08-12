package glo

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Db *gorm.DB
)

// DbConnect 连接数据库
func DbConnect() {
	db, err := gorm.Open(Config.GopsAPI.Database.Dialect, Config.GopsAPI.Database.Addr)
	if err != nil {
		log.Panicln("打开数据库失败:", err)
	}
	// 数据库日志调试
	// db.LogMode(true)

	Db = db.Set("gorm:table_options", "charset=utf8")
}

// DbDisconnect 断开数据库
func DbDisconnect() {
	Db.Close()
}
