package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(
		"mysql",
		//https://github.com/jinzhu/gorm/issues/403
		//fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Cfg.DBUsername, Cfg.DBPassword, Cfg.DBHostname, Cfg.DBPort, Cfg.DBDatabase),
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", Cfg.DBUsername, Cfg.DBPassword, Cfg.DBHostname, Cfg.DBPort, Cfg.DBDatabase),
		// fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", "root", "123456", "127.0.0.1", 3306, "app"),
	)
	if Cfg.Env == "dev" || Cfg.Env == "staging" || Cfg.Env == "test" {
		DB.LogMode(true)
	}

	if err != nil {
		panic(err)
	}

	DB.BlockGlobalUpdate(true)
}
