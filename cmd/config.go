package main

type Config struct {
	ServicePort int    `env:"SERVICE_PORT"`
	DBUser      string `env:"DB_USER"`
	DBPass      string `env:"DB_PASS"`
	DBHost      string `env:"DB_HOST"`
	DBPort      int    `env:"DB_PORT"`
	DBName      string `env:"DB_NAME"`
}
