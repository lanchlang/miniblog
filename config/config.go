package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// loadConfig loads app config,使用json格式
func LoadConfig(configFile string, config interface{}) error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}