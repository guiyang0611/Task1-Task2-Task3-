package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Config struct {
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

var (
	GormDB *gorm.DB
	SqlxDB *sqlx.DB
)

func InitDb() error {
	config, err := LoadConfig("resource/config.dev.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	// 2. 使用配置中的 DSN 初始化数据库连接
	db, err := gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	GormDB = db
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取 sql.DB 失败: %v", err)
	}
	SqlxDB = sqlx.NewDb(sqlDB, "mysql")
	log.Println("✅ Successfully connected to database!")
	return nil
}
