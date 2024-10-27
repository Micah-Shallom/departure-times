package storage

import (
	"github.com/Micah-Shallom/departure-times/utility"
	"github.com/redis/go-redis/v9"
)

var ()

type Database struct {
	Redis *redis.Client
}

var (
	DB *Database = &Database{}
	Logger *utility.Logger
)

func Connection () *Database {
	return DB
}


