package main

type GlobalConfig struct {
	UpdateInterval string `yaml:"update_interval" json:"update_interval"` // Cron expression
	DatabasePath   string `yaml:"database_path" json:"database_path"`
	OutputPath     string `yaml:"output_path" json:"output_path"`     // Directory for M3U files
	WebServerPort  string `yaml:"web_server_port" json:"web_server_port"` // Port for web interface (default: 8080)
}

type Subscription struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url" json:"url"`
	Cron string `yaml:"cron" json:"cron"` // Optional override
}

type Config struct {
	Global        GlobalConfig   `yaml:"global" json:"global"`
	Subscriptions []Subscription `yaml:"subscriptions" json:"subscriptions"`
}
