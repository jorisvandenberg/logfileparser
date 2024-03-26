package initialise

import (
	"fmt"
	"logfileparser/shared/db/connection"
	"logfileparser/shared/shared"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func EndsWithSlash(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '/'
}

func InitialiseDb() error {
	if !EndsWithSlash(shared.MyConfig.DatabasePath) {
		shared.MyConfig.DatabasePath += "/"
	}

	if _, err := os.Stat(shared.MyConfig.DatabasePath); os.IsNotExist(err) {
		err := os.MkdirAll(shared.MyConfig.DatabasePath, 0755)
		if err != nil {
			panic(err)
		}
	}

	db, err := connection.DbConnect(filepath.Join(shared.MyConfig.DatabasePath, shared.MyConfig.DatabaseName))
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS "logline" (
            "timestamp" INTEGER NOT NULL,
            "bytes_sent" INTEGER,
            "duration" INTEGER,
            "id_remote_ip" INTEGER,
            "id_remote_hostname" INTEGER,
            "id_remote_user" INTEGER,
            "id_request_method" INTEGER,
            "id_request_uri" INTEGER,
            "id_http_protocol" INTEGER,
            "id_status_code" INTEGER,
            "id_referrer" INTEGER,
            "id_user_agent" INTEGER,
            "id_cookies" INTEGER,
            "id_query_parameters" INTEGER,
            "id_headers" INTEGER,
            "id_server_name" INTEGER,
            "id_remote_logname" INTEGER,
            FOREIGN KEY("id_remote_ip") REFERENCES "remote_ip"("id"),
            FOREIGN KEY("id_remote_hostname") REFERENCES "remote_hostname"("id"),
            FOREIGN KEY("id_remote_user") REFERENCES "remote_user"("id"),
            FOREIGN KEY("id_request_method") REFERENCES "request_method"("id"),
            FOREIGN KEY("id_request_uri") REFERENCES "request_uri"("id"),
            FOREIGN KEY("id_http_protocol") REFERENCES "http_protocol"("id"),
            FOREIGN KEY("id_status_code") REFERENCES "status_code"("id"),
            FOREIGN KEY("id_referrer") REFERENCES "referrer"("id"),
            FOREIGN KEY("id_user_agent") REFERENCES "user_agent"("id"),
            FOREIGN KEY("id_cookies") REFERENCES "cookies"("id"),
            FOREIGN KEY("id_query_parameters") REFERENCES "query_parameters"("id"),
            FOREIGN KEY("id_headers") REFERENCES "headers"("id"),
            FOREIGN KEY("id_server_name") REFERENCES "server_name"("id"),
            FOREIGN KEY("id_remote_logname") REFERENCES "remote_logname"("id"),
			UNIQUE(timestamp, bytes_sent, duration, id_remote_ip, id_remote_hostname, id_remote_user, id_request_method, id_request_uri, id_http_protocol, id_status_code, id_referrer, id_user_agent, id_cookies, id_query_parameters, id_headers, id_server_name, id_remote_logname)
        );`,
		`CREATE TABLE IF NOT EXISTS "cookies" (
            "id" INTEGER NOT NULL UNIQUE,
            "cookie" TEXT NOT NULL UNIQUE,
            PRIMARY KEY("id" AUTOINCREMENT)
        );`,
		`CREATE INDEX IF NOT EXISTS idx_cookie ON cookies (cookie);`,
		`CREATE TABLE IF NOT EXISTS "finished_files" (
            "id" INTEGER NOT NULL UNIQUE,
            "hash" TEXT NOT NULL UNIQUE,
            PRIMARY KEY("id" AUTOINCREMENT)
        );`,
		`CREATE INDEX IF NOT EXISTS idx_hash ON finished_files (hash);`,
		`CREATE TABLE IF NOT EXISTS "headers" (
			"id"	INTEGER NOT NULL UNIQUE,
			"header"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_header ON headers (header);`,
		`CREATE TABLE IF NOT EXISTS "http_protocol" (
			"id"	INTEGER NOT NULL UNIQUE,
			"protocol"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_protocol ON http_protocol (protocol);`,
		`CREATE TABLE IF NOT EXISTS "query_parameters" (
			"id"	INTEGER NOT NULL UNIQUE,
			"parameters"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_parameters ON query_parameters (parameters);`,
		`CREATE TABLE IF NOT EXISTS "referrer" (
			"id"	INTEGER NOT NULL UNIQUE,
			"referrer"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_referrer ON referrer (referrer);`,
		`CREATE TABLE IF NOT EXISTS "remote_hostname" (
			"id"	INTEGER NOT NULL UNIQUE,
			"hostname"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_hostname ON remote_hostname (hostname);`,
		`CREATE TABLE IF NOT EXISTS "remote_ip" (
			"id"	INTEGER NOT NULL UNIQUE,
			"ip"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_ip ON remote_ip (ip);`,
		`CREATE TABLE IF NOT EXISTS "remote_logname" (
			"id"	INTEGER NOT NULL UNIQUE,
			"remote_logname"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_remote_logname ON remote_logname (remote_logname);`,
		`CREATE TABLE IF NOT EXISTS "remote_user" (
			"id"	INTEGER NOT NULL UNIQUE,
			"user"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_user ON remote_user (user);`,
		`CREATE TABLE IF NOT EXISTS "request_method" (
			"id"	INTEGER NOT NULL UNIQUE,
			"method"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_method ON request_method (method);`,
		`CREATE TABLE IF NOT EXISTS "request_uri" (
			"id"	INTEGER NOT NULL UNIQUE,
			"uri"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_uri ON request_uri (uri);`,
		`CREATE TABLE IF NOT EXISTS "server_name" (
			"id"	INTEGER NOT NULL UNIQUE,
			"server_name"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_server_name ON server_name (server_name);`,
		`CREATE TABLE IF NOT EXISTS "status_code" (
			"id"	INTEGER NOT NULL UNIQUE,
			"code"	INTEGER NOT NULL UNIQUE,
			"description"	TEXT,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_code ON status_code (code);`,
		`CREATE TABLE IF NOT EXISTS "user_agent" (
			"id"	INTEGER NOT NULL UNIQUE,
			"user_agent"	TEXT NOT NULL UNIQUE,
			PRIMARY KEY("id" AUTOINCREMENT)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_user_agent ON user_agent (user_agent);`,
		`INSERT OR IGNORE INTO "cookies" ("id", "cookie") VALUES (1, 'unset');`,
		`INSERT OR IGNORE INTO "headers" ("id","header") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "http_protocol" ("id","protocol") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "query_parameters" ("id","parameters") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "referrer" ("id","referrer") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "remote_hostname" ("id","hostname") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "remote_ip" ("id","ip") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "remote_logname" ("id","remote_logname") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "remote_user" ("id","user") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "request_method" ("id", "method") VALUES 
		(1, 'unset'),
			(2,'GET'),
			(3,'POST'),
			(4,'HEAD'),
			(5,'PUT'),
			(6,'DELETE'),
			(7,'OPTIONS'),
			(8,'PATCH');`,
		`INSERT OR IGNORE INTO "request_uri" ("id","uri") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "server_name" ("id","server_name") VALUES (1,'unset');`,
		`INSERT OR IGNORE INTO "status_code" ("id","code","description") VALUES (1,0,'unset'),
		(2,100,'Continue'),
		(3,101,'Switching Protocols'),
		(4,102,'Processing'),
		(5,103,'Early Hints'),
		(6,200,'OK'),
		(7,201,'Created'),
		(8,202,'Accepted'),
		(9,203,'Non-Authoritative Information'),
		(10,204,'No Content'),
		(11,205,'Reset Content'),
		(12,206,'Partial Content'),
		(13,207,'Multi-Status'),
		(14,208,'Already Reported'),
		(15,226,'IM Used'),
		(16,300,'Multiple Choices'),
		(17,301,'Moved Permanently'),
		(18,302,'Found'),
		(19,303,'See Other'),
		(20,304,'Not Modified'),
		(21,305,'Use Proxy'),
		(22,307,'Temporary Redirect'),
		(23,308,'Permanent Redirect'),
		(24,400,'Bad Request'),
		(25,401,'Unauthorized'),
		(26,402,'Payment Required'),
		(27,403,'Forbidden'),
		(28,404,'Not Found'),
		(29,405,'Method Not Allowed'),
		(30,500,'Internal Server Error'),
		(31,501,'Not Implemented'),
		(32,502,'Bad Gateway'),
		(33,503,'Service Unavailable'),
		(34,504,'Gateway Timeout'),
		(35,505,'HTTP Version Not Supported');`,
		`INSERT OR IGNORE INTO "user_agent" ("id","user_agent") VALUES (1,'unset');`,
	}

	for _, sqlStatement := range sqlStatements {
		_, err := tx.Exec(sqlStatement)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing SQL statement: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
