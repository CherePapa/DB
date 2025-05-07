package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbHost     = "localhost"
	dbPort     = 3306
	dbUser     = "root"
	dbPassword = ""
	dbName     = "pharmacy"
)

func CreateDBIfNotExists() error {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/?charset=utf8&parseTime=True", dbUser, dbPassword)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Создание БД если не существует
	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	sqlDB, _ := db.DB()
	sqlDB.Close()
	return nil
}

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	fmt.Println("Пытаюсь подключиться с DSN:", dsn) // Добавьте эту строку

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Ошибка подключения:", err) // Добавьте вывод ошибки
		return nil, err
	}
	return db, nil
}
