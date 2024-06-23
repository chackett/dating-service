package main

// Config defines application configuration, to be populated via envars
type Config struct {
	// ServicePort defines the port the web service is to be exposed on
	ServicePort int    `env:"SERVICE_PORT"`
	DBUser      string `env:"DB_USER"`
	DBPass      string `env:"DB_PASS"`
	DBHost      string `env:"DB_HOST"`
	DBPort      int    `env:"DB_PORT"`
	DBName      string `env:"DB_NAME"`
}
