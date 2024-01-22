package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

var CONFIG_FILE_NAME = "setting-local.yml"

type Config struct {
	Env Env `yaml:"env"`
}

type Env struct {
	System    string   `yaml:"system"`
	Area      string   `yaml:"area"`
	Group     string   `yaml:"group"`
	LogType   string   `yaml:"logType"`
	Host      string   `yaml:"host"`
	Level     string   `yaml:"level"`
	PrintType []string `yaml:"printType"`
	FilePath  string   `yaml:"filePath"`
	Pod       string   `yaml:"pod"`
	Gid       string   `yaml:"gid"`
}

var Cfg Config

func init() {
	loadConfig()
}

func loadConfig() {
	fnc := "loadCfgFromFile"
	var f *os.File
	var err error

	cfg := &Cfg
	f, err = os.OpenFile(CONFIG_FILE_NAME, os.O_RDONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
			f, err = os.Create(CONFIG_FILE_NAME)
			if err != nil {
				fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
				}
			}(f)

			_, err := yaml.Marshal(cfg)
			if err != nil {
				fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
			}
		} else {
			fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
		}
	} else {
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
			}
		}(f)

		data, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
		}

		err = yaml.Unmarshal(data, cfg)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s error: %v", fnc, err))
		}
	}
}
