package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SMSConfig struct {
	PhoneNumbers  []string `json:"phone_numbers"`
	SignName      string   `json:"sign_name"`
	TemplateCode  string   `json:"template_code"`
	TemplateParam string   `json:"template_param"`
}

type SingleCallByTTSConfig struct {
	CalledShowNumber string `json:"called_show_number"`
	CalledNumber     string `json:"called_number"`
	TemplateCode     string `json:"template_code"`
	TemplateParam    string `json:"template_param"`
}

type Config struct {
	AccessKeyID     string                `json:"access_key_id"`
	AccessKeySecret string                `json:"access_key_secret"`
	SMS             SMSConfig             `json:"sms"`
	SingleCallByTTS SingleCallByTTSConfig `json:"single_call_by_tts"`
}

// loadConfig loads app config.
func LoadConfig(configFile string, config *Config) error {
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