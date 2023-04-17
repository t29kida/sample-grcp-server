package config

import "os"

const (
	defaultPort = "8080"
	defaultEnv  = "development"
)

var Cfg = &Config{
	port: defaultPort,
	env:  defaultEnv,
}

type Config struct {
	port       string
	env        string
	dbUser     string
	dbPassword string
	dbName     string
	dbAddr     string
}

func (c *Config) GetDBUser() string {
	return c.dbUser
}

func (c *Config) GetDBPassword() string {
	return c.dbPassword
}

func (c *Config) GetDBName() string {
	return c.dbName
}

func (c *Config) GetDBAddr() string {
	return c.dbAddr
}

func LoadConfig() {
	// PORTを読み込む
	if port := os.Getenv("PORT"); port != "" {
		Cfg.port = port
	}

	// 環境名を読み込む
	if env := os.Getenv("ENV"); env != "" {
		Cfg.env = env
	}

	// ユーザー名を読み込む
	Cfg.dbUser = os.Getenv("DB_USER")

	// パスワードを読み込む
	Cfg.dbPassword = os.Getenv("DB_PASSWORD")

	// データベース名を読み込む
	Cfg.dbName = os.Getenv("DB_NAME")

	// データベースホスト名を読み込む
	Cfg.dbAddr = os.Getenv("DB_ADDR")
}
