package database

import (
	"database/sql"
	"log"

	"webpaint/internal/config"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	DB *sql.DB
}

func InitDB(cfg config.Config) *Database {
	mysqlCfg := mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPass,
		Net:                  "tcp",
		Addr:                 cfg.DBHost,
		DBName:               cfg.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		log.Fatal("Ошибка подключения к MySQL:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к MySQL:", err)
	}

	createTables(db)

	return &Database{DB: db}
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы пользователей:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS drawings (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		title VARCHAR(100) NOT NULL,
		data LONGTEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы рисунков:", err)
	}
}
