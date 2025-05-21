package conf

import (
	"path/filepath"
	"sync"

	"github.com/kr/pretty"
	"github.com/spf13/viper"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Server Server `yaml:"server"`
	Mongo  Mongo  `yaml:"mongo"`
	ETCD   ETCD   `yaml:"etcd"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Mongo struct {
	URI            string `yaml:"uri"`
	ConnectTimeOut int    `yaml:"connectTimeOut"`
}

type ETCD struct {
	Endpoint       []string `yaml:"endpoint"`
	ConnectTimeOut int      `yaml:"connectTimeOut"`
}

func GetConf() *Config {
	once.Do(func() {
		initConf()
	})
	return conf
}

func initConf() {
	prefix := "./conf"
	path := filepath.Join(prefix, "conf.yml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	conf = new(Config)
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	pretty.Printf("%+v\n", conf)
}
