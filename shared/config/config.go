package config

import (
	"flag"
	"fmt"
	"logfileparser/shared/shared"
	"os"

	"github.com/go-ini/ini"
)

func LoadMyConfig() {
	if err := LoadConfigFromIniFile(); err != nil {
		fmt.Println("Error loading config from ini file:", err)
	}

	ParseCommandLineFlags()
}

func LoadConfigFromIniFile() error {
	configFilePaths := []string{"config.ini", "/etc/logparser/config.ini"}

	for _, path := range configFilePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		cfg, err := ini.Load(path)
		if err != nil {
			return fmt.Errorf("failed to load config file %s: %v", path, err)
		}

		section := cfg.Section("logparser")

		if KeyExists(section, "log_filepath") {
			shared.MyConfig.LogFilePath = section.Key("log_filepath").String()
		}
		if KeyExists(section, "log_file_filter") {
			shared.MyConfig.LogFileFilter = section.Key("log_file_filter").String()
		}
		if KeyExists(section, "database_path") {
			shared.MyConfig.DatabasePath = section.Key("database_path").String()
		}
		if KeyExists(section, "database_name") {
			shared.MyConfig.DatabaseName = section.Key("database_name").String()
		}
		if KeyExists(section, "log_line_regex") {
			shared.MyConfig.LogLineRegex = section.Key("log_line_regex").String()
		}

		return nil
	}
	return fmt.Errorf("config.ini not found in any expected location")
}

func KeyExists(section *ini.Section, keyName string) bool {
	_, err := section.GetKey(keyName)
	return err == nil
}

func ParseCommandLineFlags() {
	flag.StringVar(&shared.MyConfig.LogFilePath, "logpath", shared.MyConfig.LogFilePath, "Path where Apache log files reside")
	flag.StringVar(&shared.MyConfig.LogFileFilter, "logfilefilter", shared.MyConfig.LogFileFilter, "Filter to indicate which log files to parse")
	flag.StringVar(&shared.MyConfig.DatabasePath, "dbpath", shared.MyConfig.DatabasePath, "Path where the database containing log entries will be stored")
	flag.StringVar(&shared.MyConfig.DatabaseName, "dbname", shared.MyConfig.DatabaseName, "Name of the SQLite database containing log entries")
	flag.StringVar(&shared.MyConfig.LogLineRegex, "loglineregex", shared.MyConfig.LogLineRegex, "Regular expression to parse each log line")
	flag.Parse()
}
