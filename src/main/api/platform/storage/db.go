package storage

import (
	conf "github.com/agustinrabini/items-api-project/src/main/api/config"

	"github.com/jopitnow/go-jopit-toolkit/gonosql"
)

func NewNoSQL() *gonosql.Data {
	config := getDBConfig()
	return gonosql.NewNoSQL(config)
}

func getDBConfig() gonosql.Config {
	return gonosql.Config{
		Username: conf.ConfMap.MongoUser,
		Password: conf.ConfMap.MongoPassword,
		Host:     conf.ConfMap.MongoHost,
		Database: conf.ConfMap.MongoDataBase,
	}
}
