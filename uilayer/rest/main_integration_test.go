// +build integration

package rest

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/tonyalaribe/ninja/datalayer"
	"github.com/tonyalaribe/ninja/datalayer/mongodb"
	_ "github.com/tonyalaribe/ninja/datalayer/mongodb"
)

type Config struct {
	IsProduction bool               `mapstructure:"is_production"`
	ShortName    string             `mapstructure:"short_name"`
	LongName     string             `mapstructure:"long_name"`
	DBConfig     datalayer.DBConfig `mapstructure:"db_config"`
}

func TestMain(m *testing.M) {
	var err error
	config := GetConfig()

	// dataStore is a global in main_test.go
	dataStore, err = datalayer.Connect(config.DBConfig.DriverType, config.DBConfig)
	if err != nil {
		log.Panicf("unable to connect to datastore, %v", err)
	}
	DropDB(dataStore)
	// defer DropDB(dataStore)
	os.Exit(m.Run())
}

// DropDB deletes the db in use by casting the datastore to a mongodb Datastore struct and accessing the underlying db instance.
func DropDB(db datalayer.DataStore) {
	mongoDBInstance := dataStore.(*mongodb.Datastore)
	err := mongoDBInstance.DB.DropDatabase()
	if err != nil {
		log.Panicf("unable to delete db instance with err=%v", err)
	}
}

func GetConfig() Config {
	viper.SetConfigFile("../../.ninja_test.yaml")

	var err error
	if err = viper.ReadInConfig(); err != nil {
		log.Println("Unable to read config:", err)
		os.Exit(1)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panicf("unable to decode into struct, %v", err)
	}

	if !config.IsProduction {
		log.Println("In Development Mode. Logging configuration data:")
		indentedConfig, _ := json.MarshalIndent(config, "", "\t")
		log.Printf("\n%s\n\n", indentedConfig)
	}
	return config
}
