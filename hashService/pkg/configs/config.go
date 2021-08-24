package configs

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
)

const configProvider = "consul"
const configHost = "localhost:8500"
const configAppKey = "HASH_APP_CONFIG"
const configValueType = "json"

func InitConfig() {
	err := viper.AddRemoteProvider(configProvider, configHost, configAppKey)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigType(configValueType)
}

func readConfig() {
	err := viper.ReadRemoteConfig()

	if err != nil {
		log.Fatal(err, "Can't read config")
	}
}

func Get(key string) interface{} {
	readConfig()
	return viper.Get(key)
}
