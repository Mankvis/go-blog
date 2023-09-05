package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"go-blog/pkg/logger"
	"time"
)

// DB 数据库对象
var DB *sql.DB

func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "go_blog",
		AllowNativePasswords: true,
	}

	// 准备数据库链接池
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大链接数
	DB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试链接，失败会报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSql := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
	);`
	_, err := DB.Exec(createArticlesSql)
	logger.LogError(err)
}
