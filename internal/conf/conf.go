package conf

import (
	"errors"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

type IConf interface {
	Load(configFilePath string) error
	LoggerCfg() LoggerConf
	APICfg() APIConf
}

// Conf is a struct to store params of config, which read from config.ini.
type Conf struct {

	// Params of API server
	apiServiceName string
	apiProtocol    string
	apiHost        string
	apiPort        string
	apiMode        string
	apiKeyFilePath string

	// Params of log file stored path
	infoDebugLogPath string
	warnPanicLogPath string
}

type LoggerConf struct {
	APIServiceName   string
	InfoDebugLogPath string
	WarnPanicLogPath string
}

type APIConf struct {
	APIServiceName string
	APIMode        string
	APIProtocol    string
	APIHost        string
	APIPort        string
	APIKeyFilePath string
}

// Load is used to load config.ini and set fileds of Conf
func (conf *Conf) Load(configFilePath string) error {

	// Init config reader
	confReader, err := ini.Load(configFilePath)
	if err != nil {
		return errors.New("config reader init failed: " + err.Error())
	}

	// Root path
	rootPath, err := os.Getwd()
	if err != nil {
		return errors.New("read root path failed: " + err.Error())
	}

	// Params of API server

	apiServiceName, err := conf.GetString(confReader, "API SERVER", "Service_Name")
	if err != nil {
		return errors.New("read [API SERVER] ServiceName failed: " + err.Error())
	}
	conf.apiServiceName = apiServiceName

	apiProtocol, err := conf.GetString(confReader, "API SERVER", "Protocol")
	if err != nil {
		return errors.New("read [API SERVER] Protocol failed: " + err.Error())
	}
	conf.apiProtocol = apiProtocol

	apiHost, err := conf.GetString(confReader, "API SERVER", "Host")
	if err != nil {
		return errors.New("read [API SERVER] Host failed: " + err.Error())
	}
	conf.apiHost = apiHost

	apiPort, err := conf.GetString(confReader, "API SERVER", "Port")
	if err != nil {
		return errors.New("read [API SERVER] Port failed: " + err.Error())
	}
	conf.apiPort = apiPort

	apiMode, err := conf.GetString(confReader, "API SERVER", "Mode")
	if err != nil {
		return errors.New("read [API SERVER] Mode failed: " + err.Error())
	}
	conf.apiMode = apiMode

	apiKeyFilePath, err := conf.GetString(confReader, "API SERVER", "APIKey_File_Path")
	if err != nil {
		return errors.New("read [API SERVER] APIKey_File_Path failed: " + err.Error())
	}
	conf.apiKeyFilePath = path.Join(rootPath, apiKeyFilePath)

	// Params of log file stored path

	infoDebugLogPath, err := conf.GetString(confReader, "FILE STORED PATH", "Info_Debug_Log_Path")
	if err != nil {
		return errors.New("read [FILE STORED PATH] Info_Debug_Log_Path failed: " + err.Error())
	}
	conf.infoDebugLogPath = path.Join(rootPath, infoDebugLogPath)

	warnPanicLogPath, err := conf.GetString(confReader, "FILE STORED PATH", "Warn_Panic_Log_Path")
	if err != nil {
		return errors.New("read [FILE STORED PATH] Warn_Panic_Log_Path failed: " + err.Error())
	}
	conf.warnPanicLogPath = path.Join(rootPath, warnPanicLogPath)

	return nil
}

func (conf *Conf) LoggerCfg() LoggerConf {
	loggerConf := LoggerConf{
		APIServiceName:   conf.apiServiceName,
		InfoDebugLogPath: conf.infoDebugLogPath,
		WarnPanicLogPath: conf.warnPanicLogPath,
	}
	return loggerConf
}

func (conf *Conf) APICfg() APIConf {
	apiConf := APIConf{
		APIServiceName: conf.apiServiceName,
		APIMode:        conf.apiMode,
		APIProtocol:    conf.apiProtocol,
		APIHost:        conf.apiHost,
		APIPort:        conf.apiPort,
		APIKeyFilePath: conf.apiKeyFilePath,
	}
	return apiConf
}

// GetString read string from section with key
func (conf *Conf) GetString(confReader *ini.File, section string, key string) (string, error) {
	if confReader == nil {
		return "", errors.New("no conf reader")
	}

	s := confReader.Section(section)
	if s == nil {
		return "", errors.New("no such section")
	}

	value := s.Key(key).String()
	if value == "" {
		return "", errors.New("no such key")
	}

	return value, nil
}

// GetInt read int from section with key
func (conf *Conf) GetInt(confReader *ini.File, section string, key string) (int, error) {
	if confReader == nil {
		return 0, errors.New("no conf reader")
	}

	s := confReader.Section(section)
	if s == nil {
		return 0, errors.New("no such section")
	}

	valueInt, _ := s.Key(key).Int()

	if valueInt == 0 {
		return 0, errors.New("no such key")
	}

	return valueInt, nil
}
