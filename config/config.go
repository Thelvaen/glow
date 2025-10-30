package config

import (
	"log"

	"github.com/spf13/viper"
)

type Program struct {
	Label string            `mapstructure:"label"`
	Cmd   string            `mapstructure:"cmd"`
	Cwd   string            `mapstructure:"cwd"`
	Env   map[string]string `mapstructure:"env"`
}

type Button struct {
	Label    string    `mapstructure:"label"`
	Programs []Program `mapstructure:"programs"`
}

type Config struct {
	Title   string   `mapstructure:"title"`
	Buttons []Button `mapstructure:"buttons"`
}

func Load(configPath string) Config {
	v := viper.New()
	v.SetConfigType("yaml")

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	v.SetDefault("title", "Quick Launcher")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return cfg
}
