package main

import (
	"encoding/json"
	"fmt"
	"os"
	"savemylink/database"
)

type Config struct {
	Port     int                `json:"port"`
	Env      string             `json: "env"`
	Database *database.DbConfig `json: "database:`
}

func LoadConfig(isProd bool) *Config {
	if !isProd {
		fmt.Println("Successfully loaded dev config")
		return devConfig()
	}

	f, err := os.Open(".config")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)
	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully loaded prod config")
	return &c
}

func devConfig() *Config {
	return &Config{
		Port:     3000,
		Env:      "dev",
		Database: database.DevDbConfig(),
	}
}
