package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("config/.")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Config/Server: error read config yaml : %v", err)
	}
	var config Config
	if err := viper.UnmarshalKey("server", &config); err != nil {
		logrus.Fatalf("Config/Server: error unmarshal config yaml : %v", err)
	}

	return &config
}

type ConfigPQ struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Ssl      string `yaml:"ssl"`
}

func NewConfigPQ() *ConfigPQ {
	viper.SetConfigName("config")
	viper.AddConfigPath("/home/midas/GolandProjects/TTMusic/config/")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Config/PQ: error read config yaml: %v", err)
	}
	var configPQ ConfigPQ
	if err := viper.UnmarshalKey("postgresql", &configPQ); err != nil {
		logrus.Fatalf("Config/PQ: error unmarshal config yaml : %v", err)
	}
	return &configPQ
}

type KeyGenius struct {
	Key string
}
type KeyOpenAI struct {
	Key string
}
type KeyYouTube struct {
	Key string
}

func (k *KeyGenius) GetAPIKey() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/home/midas/GolandProjects/TTMusic/config/")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("Config/Genius: error read config yaml : %v", err)
	}
	if err := viper.UnmarshalKey("key_genius", &k.Key); err != nil {
		logrus.Panicf("Server/MS_API: error unmarshal config yaml : %v", err)
	}
}
func (k *KeyOpenAI) GetAPIKey() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/home/midas/GolandProjects/TTMusic/config/")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("Config/OpenAI: error read config yaml : %v", err)
	}
	if err := viper.UnmarshalKey("key_openai", &k.Key); err != nil {
		logrus.Panicf("Server/MS_API: error unmarshal config yaml : %v", err)
	}
}
func (k *KeyYouTube) GetAPIKey() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/home/midas/GolandProjects/TTMusic/config/")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("Config/OpenAI: error read config yaml : %v", err)
	}
	if err := viper.UnmarshalKey("key_youtube", &k.Key); err != nil {
		logrus.Panicf("Server/MS_API: error unmarshal config yaml : %v", err)
	}
}
