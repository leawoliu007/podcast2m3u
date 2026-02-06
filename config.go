package main

type GlobalConfig struct {
	UpdateInterval string `yaml:"update_interval"` // Cron expression, e.g., "0 * * * *" for every hour
	DatabasePath   string `yaml:"database_path"`
}

type Subscription struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	Cron       string `yaml:"cron"` // Optional override
	OutputPath string `yaml:"output_path"`
}

type Config struct {
	Global        GlobalConfig   `yaml:"global"`
	Subscriptions []Subscription `yaml:"subscriptions"`
}
