package repository

import (
	"authService/config"
	"authService/db"
	"sync"
)


var (
	initOnce sync.Once
)


func Init() {
	initOnce.Do(func() {
		db.InitDB(config.GetDBConfig())
	})
}