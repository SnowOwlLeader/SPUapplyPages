package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	OAuth    OAuthConfig    `yaml:"oauth"`
	Database DatabaseConfig `yaml:"database"`
	Google   struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RefreshToken string `yaml:"refresh_token"`
	} `yaml:"google"`
}

type ServerConfig struct {
	Port string `yaml:"port" default:"8080"`
}

type OAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
	TokenURL     string `yaml:"token_url"`
	AuthURL      string `yaml:"auth_url"`
	UserInfoURL  string `yaml:"userinfo_url"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

var GlobalConfig Config

func Init() error {
	// 获取当前文件所在目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	// 读取配置文件
	configPath := filepath.Join(dir, "config", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		// 尝试读取相对路径
		data, err = os.ReadFile("config/config.yaml")
		if err != nil {
			return err
		}
	}

	// 解析YAML
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return err
	}

	return nil
}
