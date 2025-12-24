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

// Action is a step in a scenario: start cmds in a tab or wait for N seconds
type Action struct {
	Tab  string   `mapstructure:"tab"`
	Cmds []string `mapstructure:"cmds"`
	Wait int      `mapstructure:"wait"`
}

type Scenario struct {
	Name    string   `mapstructure:"name"`
	Actions []Action `mapstructure:"actions"`
}

type Config struct {
	Title     string     `mapstructure:"title"`
	Buttons   []Button   `mapstructure:"buttons"`
	Scenarios []Scenario `mapstructure:"scenarios"`
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

	// Backwards compatibility: convert Buttons -> Scenarios
	if len(cfg.Scenarios) == 0 && len(cfg.Buttons) > 0 {
		for _, b := range cfg.Buttons {
			var s Scenario
			s.Name = b.Label
			for _, p := range b.Programs {
				act := Action{Tab: p.Label, Cmds: []string{p.Cmd}}
				s.Actions = append(s.Actions, act)
			}
			cfg.Scenarios = append(cfg.Scenarios, s)
		}
	}

	return cfg
}
