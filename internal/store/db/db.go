package db

import (
	"os"

	"github.com/s30899-pj/HomePiggyBank_byt2025-26_52c/internal/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func open(dbName string) (*gorm.DB, error) {
	err := os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(dbName), &gorm.Config{})
}

func MustOpen(dbName string) *gorm.DB {
	if dbName == "" {
		dbName = "hpb.db"
	}

	db, err := open(dbName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&store.User{}, &store.Session{}, &store.Household{}, &store.Membership{}, &store.Expense{}, &store.ExpenseShare{}, &store.Report{})
	if err != nil {
		panic(err)
	}

	return db
}
