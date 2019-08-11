package glo

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// DatabaseConfig Struct
type DatabaseConfig struct {
	Dialect     string `yaml:"dialect"`
	Addr        string `yaml:"addr"`
	AutoMigrate bool   `yaml:"automigrate"`
}

// RedisConfig Struct
type RedisConfig struct {
	MaxIdle        int    `yaml:"max_idle"`
	IdleTimeoutSec int    `yaml:"idle_timeout_sec"`
	DefaultExpried int    `yaml:"expried_sec"`
	Addr           string `yaml:"addr"`
	Password       string `yaml:"password"`
	LoopSec        int    `yaml:"loop_sec"`
	Expried        int    `yaml:"expried"`
}

type LogConfig struct {
	Mode string `yaml:"mode"`
	Path string `yaml:"path"`
}

type DingDingConfig struct {
	API string `yaml:"api"`
}
type NotifyConfig struct {
	DingDing DingDingConfig `yaml:"dingding"`
}
type AlertCfg struct {
	AlertChan int64 `yaml:"alert_chan"`
}

type ModifyConfig struct {
	Details            string `yaml:"details"`
	SLAProcessAlert    int    `yaml:"sla_process_alert"`
	CloudModifyDetails string `yaml:"cloud_modify_url"`
}

type TicketConfig struct {
	Details            string `yaml:"details"`
	SLAProcessAlert    int    `yaml:"sla_process_alert"`
	CloudTicketDetails string `yaml:"cloud_ticket_url"`
}

// GopsAPICfg Struct
type GopsAPICfg struct {
	EncryptKey string         `yaml:"encrypt_key"`
	Enable     bool           `yaml:"enable"`
	MaxRequest string         `yaml:"max_request"`
	Deadline   int            `yaml:"deadline"`
	Database   DatabaseConfig `yaml:"database"`
	Redis      RedisConfig    `yaml:"redis"`
	ServerPort int            `yaml:"server_port"`
	Log        LogConfig      `yaml:"log"`
	Notify     NotifyConfig   `yaml:"notify"`
	Ticket     TicketConfig   `yaml:"ticket"`
	Modify     ModifyConfig   `yaml:"modify"`
	Alert      AlertCfg       `yaml:"alert"`
}

// GlobalConfig Struct
type GlobalConfig struct {
	GopsAPI GopsAPICfg `yaml:"gops_api"`
}

var (
	// Config Params
	Config GlobalConfig
	// ConfigFile FileName
	ConfigFile string
)

// ParseConfig Function
func ParseConfig(file string) (err error) {
	ConfigFile = file
	var c GlobalConfig
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return
	}
	Config = c
	return
}
