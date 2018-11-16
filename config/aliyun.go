package config

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

