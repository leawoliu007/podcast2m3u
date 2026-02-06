package main

type GlobalConfig struct {
	UpdateInterval string `yaml:"update_interval"` // Cron expression
	DatabasePath   string `yaml:"database_path"`
	OutputPath     string `yaml:"output_path"`     // Directory for M3U files
}

type Subscription struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
	Cron string `yaml:"cron"` // Optional override
}

type Config struct {
	Global        GlobalConfig   `yaml:"global"`
	Subscriptions []Subscription `yaml:"subscriptions"`
}
