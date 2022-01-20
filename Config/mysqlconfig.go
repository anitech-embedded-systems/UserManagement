package dbconfig

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type MysqlConfig struct {
	// Connection struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Table    string `yaml:"table"`
	Database string `yaml:"database"`
	// }
}

type Config struct {
	Config  MysqlConfig `yaml:"mysql"`
	Default string      `yaml:"default"`
}

func Get() *Config {
	var sqlconf Config
	yamlFile, err := ioutil.ReadFile("config-stage.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &sqlconf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &sqlconf
}
