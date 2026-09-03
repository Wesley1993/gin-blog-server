package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	App           AppConfig       `mapstructure:"app"`
	Database      DatabaseConfig  `mapstructure:"database"`
	Redis         RedisConfig     `mapstructure:"redis"`
	JWT           JWTConfig       `mapstructure:"jwt"`
	Elasticsearch ESConfig        `mapstructure:"elasticsearch"`
	Log           LogConfig       `mapstructure:"log"`
	Migration     MigrationConfig `mapstructure:"migration"`
	Swagger       SwaggerConfig   `mapstructure:"swagger"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // 小时
	Prefix string `mapstructure:"prefix"`
}

// ESConfig Elasticsearch配置
type ESConfig struct {
	Addresses   []string `mapstructure:"addresses"`
	IndexPrefix string   `mapstructure:"index_prefix"`
	Username    string   `mapstructure:"username"`
	Password    string   `mapstructure:"password"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// MigrationConfig 数据库迁移配置
type MigrationConfig struct {
	Dir string `mapstructure:"dir"`
}

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	BasePath string `mapstructure:"base_path"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")

	// 支持环境变量覆盖配置，如 DATABASE_HOST 覆盖 database.host
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}
