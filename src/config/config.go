package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)

type TwitterConfig struct {
	Consumerkey   string `yaml:"consumerKey"`
	ConumerSecret string `yaml:"consumerSecret"`
	Token         string `yaml:"token"`
	TokenSecret   string `yaml:"tokenSecret"`
}
type TelegramConfig struct {
	ApiKey string `yaml:"apiKey"`
	UserId int    `yaml:"userId"`
}
type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"dbName"`
}
type Config struct {
	Twitter  TwitterConfig  `yaml:"twitter"`
	Telegram TelegramConfig `yaml:"telegram"`
	Mysql    MysqlConfig    `yaml:"mysql"`
}

func (config *Config) GetDbConnectionString() string {
	return config.Mysql.Username + ":" + config.Mysql.Password + "@tcp(" + config.Mysql.Host + ":" + config.Mysql.Port + ")/" + config.Mysql.DbName + "?parseTime=true"
}

func LoadConfig(configFileName string) (config *Config, err error) {
	f, err := os.Open(configFileName)
	config = new(Config)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(byteValue, config)
	if err != nil {
		return nil, err
	}
	log.Println("Config file correctly lodaed.")
	return
}
