package configs

import (
    "sync"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    db   *gorm.DB
    once sync.Once
)

func GetDB() *gorm.DB {
    once.Do(func() {
        dsn := "host=localhost user=postgres password=postgres dbname=todo_db port=5432 sslmode=disable"
        var err error
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            panic("Failed to connect to database")
        }
    })
    return db
}