package database

import (
	"savemylink/save"
	"savemylink/services"
)

type Db interface {
	GetUserDb() services.UserDB
	GetArticleDb() save.ArticleDB
	Close()
}

type DbConfig struct {
	Host     string `json:"host"` // JSON tags
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name`
}

func DevDbConfig() *DbConfig {
	return &DbConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "dev_group",
		Password: "Ab123456!",
		DbName:   "savemylink",
	}
}

func NewDbService(config *DbConfig) Db {
	if config.Host != "localhost" {
		return GetDiskDb(config)
	}

	return NewMemoryDb()
}
