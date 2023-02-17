package lib

import (
	"log"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerHost   string `mapstructure:"SERVER_HOST"`
	ServerPort   string `mapstructure:"SERVER_PORT"`
	SwaggerPort  string `mapstructure:"SWAGGER_PORT"`
	Environment  string `mapstructure:"ENV"`
	LogOutput    string `mapstructure:"LOG_OUTPUT"`
	DBUsername   string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASS"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBName       string `mapstructure:"DB_NAME"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	DBCollection string `mapstructure:"DB_COLLECTION"`
	DBTrash      string `mapstructure:"DB_TRASH_COLLECTION"`
	DbUrl        string `mapstructure:"MONGO_URL"`
}

// NewEnv creates a new environment
func NewEnv() Env {

	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	return env
}
