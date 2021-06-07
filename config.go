package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const DEFAULT_CONFIG_FILE_NAME string = "/opt/tinyurl/config.yaml"
const DEFAULT_DB_FILE_NAME string = "/opt/tinyurl/tinyurl.db"
const DEFAULT_LOG_FILE_NAME string = "/opt/tinyurl/tinyurl.log"
const DEFAULT_LOG_OUTPUT_MODE string = "stderr"
const DEFAULT_LOG_LEVEL string = "info"
const DEFAULT_HTTP_PORT int = 80
const DEFAULT_PROTOCOL string = "http"

type Config struct {
	DBFileName    string `yaml:"DBFileName"`
	LogFileName   string `yaml:"LogFileName"`
	LogOutputMode string `yaml:"LogOutputMode"`
	LogLevel      string `yaml:"LogLevel"`
	HTTPPort      int    `yaml:"HTTPPort"`
	Protocol      string `yaml:"Protocol"`
}

func NewConfig(fileName string) (*Config, error) {
	var fName string
	fName = fileName
	if fName == "" {
		fmt.Fprintf(os.Stderr, "Config file name is empty. So \"%s\" will be read as default.\n", DEFAULT_CONFIG_FILE_NAME)
		fName = DEFAULT_CONFIG_FILE_NAME
	}

	var cfg Config

	// if config file is not found.
	if _, err := os.Stat(fName); err != nil {
		if fileName != "" {
			fmt.Fprintf(os.Stderr, "ConfigError: Specified config file \"%s\" was not found.\n", fName)
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "ConfigError: Default config file \"%s\" was not found. So, application will be started with default setting.\n", fName)
		return createDefaultConfig(), nil
	}

	// read config file.
	buf, err := ioutil.ReadFile(fName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadConfigError: Config file \"%s\" couldn't be read.\n", fName)
		return nil, err
	}
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "MarshalConfigError: Config file \"%s\" couldn't be parsed as yaml file.\n", fName)
	}
	// initialize params not specified value to.
	if cfg.DBFileName == "" {
		cfg.DBFileName = DEFAULT_DB_FILE_NAME
	}
	if cfg.LogFileName == "" {
		cfg.LogFileName = DEFAULT_LOG_FILE_NAME
	}
	if cfg.LogOutputMode == "" {
		cfg.LogOutputMode = DEFAULT_LOG_OUTPUT_MODE
	} else {
		if _, is := OUTPUT_MODE[cfg.LogOutputMode]; !is {
			return nil, errors.New(fmt.Sprintf("Log output mode '%d' is invalid (valid: 0-2)\n", cfg.LogOutputMode))
		}
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = DEFAULT_LOG_LEVEL
	} else {
		if _, is := LOG_LEVEL[cfg.LogLevel]; !is {
			return nil, errors.New(fmt.Sprintf("Log level '%s' is invalid (valid: debug,info,warn,error)\n", cfg.LogLevel))
		}
	}
	if cfg.HTTPPort == 0 {
		cfg.HTTPPort = DEFAULT_HTTP_PORT
	}
	if cfg.Protocol == "" {
		cfg.Protocol = DEFAULT_PROTOCOL
	} else {
		if !(cfg.Protocol == "http" || cfg.Protocol == "https") {
			return nil, errors.New(fmt.Sprintf("Protocol '%s' is invalid (valid: http,https)\n", cfg.Protocol))
		}
	}

	return &cfg, err
}

func createDefaultConfig() *Config {
	return &Config{
		DBFileName:    DEFAULT_DB_FILE_NAME,
		LogFileName:   DEFAULT_LOG_FILE_NAME,
		LogOutputMode: DEFAULT_LOG_OUTPUT_MODE,
		LogLevel:      DEFAULT_LOG_LEVEL,
		HTTPPort:      DEFAULT_HTTP_PORT,
		Protocol:      DEFAULT_PROTOCOL,
	}
}
