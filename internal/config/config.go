package config

import "github.com/spf13/viper"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
type JWTConfig struct {
	Secret string
	Issuer string
}
type Config struct {
	Listen string
	JWT    *JWTConfig
	DB     *DBConfig
}

var Cfg Config

func LoadConfig() {
	//viper.SetConfigFile("./config/config.toml")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".internal")
	viper.AddConfigPath("./internal/config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(err)
	}
}
