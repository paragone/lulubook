package utils

import (
	"errors"
	"os/exec"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)

type Config struct{
	Port 			int	//server port
	LogPath			string	//Logfile path
	LogLevel        	int
	MongodbServer		string
}

var CONFIG = new(Config)

func init(){
	var configFilePath string
	//path := getCurrentPath()
	path := os.Getenv("GOPATH")
	configFilePath = filepath.Join(path, "src/main/config.json")

	if configFilePath != "" {
		cfgData, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			panic("Failed to open found cfgFile " + configFilePath)
		}

		err = json.Unmarshal(cfgData, &CONFIG)
		if err != nil {
			panic("Read config file failed:" + err.Error())
		}
	}
	//然后再校验
	err := CONFIG.verify()
	if err != nil {
		panic("Read config file failed:" + err.Error())

	}
}

func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	i := strings.LastIndex(s, "\\")
	path := string(s[0 : i+1])
	return path
}

func (cfg *Config) verify() error {
	if cfg == nil {
		return  errors.New("cfg nil")
	}

	if cfg.Port == 0 {
		return errors.New("Required field 'Port' not found in config")
	}

	if cfg.MongodbServer == "" {
		return errors.New("Required field 'MongodbServer' not found in config")
	}

	return nil
}