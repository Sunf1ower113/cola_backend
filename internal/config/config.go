package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Listener struct {
		Protocol     string `json:"protocol"`
		Host         string `json:"host"`
		Port         string `json:"port"`
		IdleTimeout  int    `json:"idle_timeout"`
		WriteTimeout int    `json:"write_timeout"`
		ReadTimeout  int    `json:"read_timeout"`
	} `json:"listener"`
	Database struct {
		DbDriver string `json:"db_driver"`
		DbName   string `json:"db_name"`
		Config   string `json:"config"`
	} `json:"storage"`
}

func LoadConfiguration(file string) (cfg *Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panic(err)
		}
	}(f)
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return
	}
	return
}
