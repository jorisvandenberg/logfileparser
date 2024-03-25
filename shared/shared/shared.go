package shared

import (
	"time"
)

type LogLine struct {
	Remote_ip        string
	Remote_hostame   string
	Remote_user      string
	Timestamp        time.Time
	Request_method   string
	Request_URI      string
	HTTP_Protocol    string
	Status_code      int
	Bytes_sent       int
	Referrer         string
	User_Agent       string
	Duration         time.Duration
	Cookies          string
	Query_Parameters string
	Headers          string
	Server_name      string
	Remote_logname   string
}

type Config struct {
	LogFilePath   string
	LogFileFilter string
	DatabasePath  string
	DatabaseName  string
	LogLineRegex  string
}

var (
	MyConfig = Config{
		LogFilePath:   "/var/log/apache2",
		LogFileFilter: "access",
		DatabasePath:  "/tmp/",
		DatabaseName:  "apachelog.db",
		LogLineRegex:  `^(?P<RemoteIP>\S+) (?P<RemoteLogname>\S+) (?P<AuthUser>\S+) \[(?P<Timestamp>[^\]]+)\] "(?P<RequestMethod>\S+)(?: +(?P<RequestURI>[^\s"]*)[^"]*)?" (?P<StatusCode>\d+) (?P<BytesSent>\d+) "(?P<Agent>[^\"]*)"`,
	}
)
