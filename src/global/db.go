package global

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *sql.DB

// SetDatabase 初始化数据库连接
func SetDatabase() error {
	config := Config.Section("database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&allowNativePasswords=true",
		config.Key("user").String(),
		config.Key("password").String(),
		config.Key("ip").String(),
		config.Key("port").String(),
		config.Key("dbName").String(),
		config.Key("charset").String())
	var err error
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
	}

	Db.SetConnMaxLifetime(3600 * time.Second)
	Db.SetMaxIdleConns(20) //SetMaxIdleConns用于设置闲置的连接数
	Db.SetMaxOpenConns(20) //SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制

	if err := Db.Ping(); err != nil {
		log.Println(err)
		//error
		fmt.Println("Connect to mysql error")
	}
	return nil
}

var Dbs *gorm.DB

func SetDatabaseS() {
	conf := Config.Section("database")
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.Key("user").String(),
		conf.Key("password").String(),
		conf.Key("ip").String(),
		conf.Key("port").String(),
		conf.Key("dbName").String(),
	)
	db, err := gorm.Open("mysql", connArgs)
	if err != nil {
		log.Println("创建数据库连接失败:%v", err)
		return
	}
	db.DB().SetConnMaxLifetime(3600 * time.Second)
	db.DB().SetMaxIdleConns(20) //SetMaxIdleConns用于设置闲置的连接数
	db.DB().SetMaxOpenConns(20) //SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制
	Dbs = db
}
