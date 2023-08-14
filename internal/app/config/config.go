package config

import (
	"time"
)

// configure param
type Config struct {
	Global   Global   `yaml:"Global"`
	HTTP     HTTP     `yaml:"HTTP"`
	MinIO    MinIO    `yaml:"MinIO"`
	Security Security `yaml:"Security"`
	MySQL    MySQL
	Log      Log `yaml:"Log"`
}

type Global struct {
	Debug bool `yaml:"Debug"`
}
type HTTP struct {
	Host         string        `yaml:"Host"`
	Port         int           `yaml:"Port"`
	ReadTimeout  time.Duration `yaml:"ReadTimeout"`
	WriteTimeout time.Duration `yaml:"WriteTimeout"`
	IdleTimeout  time.Duration `yaml:"IdleTimeout"`
}

type MinIO struct {
	Endpoint        string          `yaml:"Endpoint"`
	AccessKeyID     string          `yaml:"AccesskeyID"`
	SecretAccessKey string          `yaml:"SecretAccessKey"`
	UseSSL          bool            `yaml:"UseSSL"`
	FileExpired     []time.Duration `yaml:"FileExpired"`
	UserExpired     time.Duration   `yaml:"UserExpired"`
	AuthParam       time.Duration   `yaml:"AuthParam"`
}

type MySQL struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DB       string `yaml:"DB"`
}
type Security struct {
	SecureKey      string        `yaml:"SecureKey"`
	TokenSecretKey string        `yaml:"TokenSecretKey"`
	TokenExpired   time.Duration `yaml:"TokenExpired"`
}

type Log struct {
	Level     int  `yaml:"Level"`
	AccessLog bool `yaml:"AccessLog"`
}
