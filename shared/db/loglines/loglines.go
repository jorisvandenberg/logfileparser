package loglines

import (
	"database/sql"
	"errors"
	"fmt"
	"logfileparser/shared/db/connection"
	"logfileparser/shared/shared"

	_ "github.com/mattn/go-sqlite3"
)

func GetIDS(db *sql.DB, isstring bool, input interface{}, tablename string, columname string) int {
	var id int
	var isemtpy bool
	if isstring {
		if input == "" {
			isemtpy = true
		} else {
			isemtpy = false
		}
	} else {
		if input == 0 {
			isemtpy = true
		} else {
			isemtpy = false
		}
	}
	if !isemtpy {
		//element not empty, see if cookie already exists in the cookies table
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM "+tablename+" WHERE "+columname+" = ?", input).Scan(&count)
		if err != nil {
			return 1
		}
		//if count > 0, select the id
		if count > 0 {
			err := db.QueryRow("SELECT id FROM "+tablename+" WHERE "+columname+" = ?", input).Scan(&id)
			if err != nil {
				return 1
			}
		} else {
			//insert a new cookie
			_, err = db.Exec("INSERT INTO "+tablename+" ("+columname+") VALUES (?)", input)
			if err != nil {
				return 1
			}
			err := db.QueryRow("SELECT id FROM "+tablename+" WHERE "+columname+" > ?", input).Scan(&id)
			if err != nil {
				return 1
			}
		}
	} else {
		//struct element empty, going for a not set value with id 1
		id = 1
	}

	return id
}

func AddLogLine(MyLogLine shared.LogLine) error {
	db, err := connection.DbConnect(shared.MyConfig.DatabasePath + shared.MyConfig.DatabaseName)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return errors.New(fmt.Sprintf("Error creating connection: %s", err))
	}
	defer db.Close()

	//first get all potential id's of all tables with referential integrity
	id_cookies := GetIDS(db, true, MyLogLine.Cookies, "cookies", "cookie")
	id_headers := GetIDS(db, true, MyLogLine.Headers, "headers", "header")
	id_http_protocol := GetIDS(db, true, MyLogLine.HTTP_Protocol, "http_protocol", "protocol")
	id_query_parameters := GetIDS(db, true, MyLogLine.Query_Parameters, "query_parameters", "parameters")
	id_referrer := GetIDS(db, true, MyLogLine.Referrer, "referrer", "referrer")
	id_remote_hostname := GetIDS(db, true, MyLogLine.Remote_hostame, "remote_hostname", "hostname")
	id_remote_ip := GetIDS(db, true, MyLogLine.Remote_ip, "remote_ip", "ip")
	id_remote_logname := GetIDS(db, true, MyLogLine.Remote_logname, "remote_logname", "remote_logname")
	id_remote_user := GetIDS(db, true, MyLogLine.Remote_user, "remote_user", "user")
	id_request_method := GetIDS(db, true, MyLogLine.Request_method, "request_method", "method")
	id_request_uri := GetIDS(db, true, MyLogLine.Request_URI, "request_uri", "uri")
	id_server_name := GetIDS(db, true, MyLogLine.Server_name, "server_name", "server_name")
	id_user_agent := GetIDS(db, true, MyLogLine.User_Agent, "user_agent", "user_agent")
	id_status_code := GetIDS(db, false, MyLogLine.Status_code, "status_code", "code")
	//fmt.Printf("%d %d %d %d %d %d %d %d %d %d %d %d %d %d\n", id_cookies, id_headers, id_http_protocol, id_query_parameters, id_referrer, id_remote_hostname, id_remote_ip, id_remote_logname, id_remote_user, id_request_method, id_request_uri, id_server_name, id_user_agent, id_status_code)
	//now add logline if it doesn't already exist

	myTimestamp := MyLogLine.Timestamp
	unixTimestamp := myTimestamp.Unix()
	myDuration := MyLogLine.Duration
	durationInSeconds := int(myDuration.Seconds())

	_, err = db.Exec("INSERT OR IGNORE INTO logline (timestamp, bytes_sent, duration, id_remote_ip, id_remote_hostname, id_remote_user, id_request_method, id_request_uri, id_http_protocol, id_status_code, id_referrer, id_user_agent, id_cookies, id_query_parameters, id_headers, id_server_name, id_remote_logname) VALUES (?, ? , ? , ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", unixTimestamp, MyLogLine.Bytes_sent, durationInSeconds, id_remote_ip, id_remote_hostname, id_remote_user, id_request_method, id_request_uri, id_http_protocol, id_status_code, id_referrer, id_user_agent, id_cookies, id_query_parameters, id_headers, id_server_name, id_remote_logname)
	if err != nil {
		return err
	}

	return nil
}
