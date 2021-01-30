package library

import (
	"fmt"

	"github.com/kecci/goscription/models"
	"github.com/ory/viper"
)

var (
	// ConfigFile name
	configFile = "config/config.toml"
	// ConfigType name
	configType = "toml"
)

// NewConfig init config
func NewConfig() models.Config {
	conf := &models.Config{}
	err := viper.Unmarshal(conf)

	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return *conf
}

// InitConfig initialize config to viper
func InitConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}
