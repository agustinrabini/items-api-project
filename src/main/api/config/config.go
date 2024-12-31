package config

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LogRatio     = 100
	LogBodyRatio = 100
)

var (
	InternalBasePricesClient = "http://prices:8080"
	InternalBaseShopsClient  = "http://shops:8080"
)

// Configuration structure
type Configuration struct {
	APIRestServerHost string `mapstructure:"api_host"`
	APIRestServerPort string `mapstructure:"api_port"`
	APIRestUsername   string `mapstructure:"api_username"`
	APIRestPassword   string `mapstructure:"api_password"`
	APIBaseEndpoint   string `mapstructure:"api_base_endpoint"`
	LoggingPath       string `mapstructure:"api_logpath"`
	LoggingFile       string `mapstructure:"api_logfile"`
	LoggingLevel      string `mapstructure:"api_loglevel"`
	MongoUser         string `mapstructure:"MONGO_USERNAME"`
	MongoPassword     string `mapstructure:"MONGO_PASSWORD"`
	MongoHost         string `mapstructure:"MONGO_HOST"`
	MongoDataBase     string `mapstructure:"MONGO_DATABASE"`
}

// ConfMap Config is package struct containing conf params
var ConfMap Configuration

func Load() {
	// Setting defaults if the config not read
	// API
	viper.SetDefault("api_host", "127.0.0.1")
	viper.SetDefault("api_port", ":8080")
	viper.SetDefault("api_username", "agus")
	viper.SetDefault("api_password", "changeme")
	viper.SetDefault("api_base_endpoint", "http://nginx")

	// LOG
	viper.SetDefault("api_logpath", "/var/log")
	viper.SetDefault("api_logfile", "api.log")
	viper.SetDefault("api_loglevel", "trace")

	viper.SetDefault("MONGO_USERNAME", "")
	viper.SetDefault("MONGO_PASSWORD", "")
	viper.SetDefault("MONGO_HOST", "")
	viper.SetDefault("MONGO_DATABASE", "")

	// Read the config file
	viper.AutomaticEnv()

	err := viper.Unmarshal(&ConfMap)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %+v", err)
	}

	spew.Dump(ConfMap)

	fmt.Println("\n All good!!")
}
