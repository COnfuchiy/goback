package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

const EnvDebug = "DEBUG"
const EnvDevelopment = "DEV"
const EnvDebugFilename = ".env.debug"
const EnvDefaultFilename = ".env"
const EnvDevelopmentFilename = ".env.dev"

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ServerPort             string `mapstructure:"SERVER_PORT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int64  `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int64  `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	EntitiesPerPage        int64  `mapstructure:"ENTITIES_PER_PAGE"`
}

func NewEnv() *Env {
	env := Env{}
	envType := os.Getenv("GO_ENV")
	var envFilename string

	switch envType {
	case EnvDebug:
		envFilename = EnvDebugFilename
		break
	case EnvDevelopment:
		envFilename = EnvDevelopmentFilename
		break
	default:
		envFilename = EnvDefaultFilename
	}

	viper.SetConfigFile(envFilename)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Can't find the file %s: ", envFilename))
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Environment can't be loaded of file %s: ", envFilename))
	}

	return &env
}
