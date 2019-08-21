package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	if err := c.initConfig(); err != nil {
		return err
	}

	c.watchConfig()
	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 解析配置文件
	} else {
		viper.AddConfigPath("conf") // 默认配置文件路径
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml") // 配置文件的格式
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("FRUIT") // 设置环境变量的前缀为FRUIT
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// viper 解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// 监控配置文件并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s\n", e.Name)
	})
}
