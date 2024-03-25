package parser

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"logfileparser/shared/db/connection"
	"logfileparser/shared/db/loglines"
	"logfileparser/shared/shared"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	//"errors"
)

func CheckifFileHandled(asciichecksum string) bool {
	db, err := connection.DbConnect(shared.MyConfig.DatabasePath + shared.MyConfig.DatabaseName)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return false
	}
	defer db.Close()
	//finished_files(hash)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM finished_files WHERE hash = ?", asciichecksum).Scan(&count)
	if err != nil {
		return true
	}
	if count > 0 {
		return false
	}
	return true
}

func SetFileHandled(asciichecksum string) bool {
	db, err := connection.DbConnect(shared.MyConfig.DatabasePath + shared.MyConfig.DatabaseName)
	if err != nil {
		fmt.Println("Error creating connection:", err)
		return false
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM finished_files WHERE hash = ?", asciichecksum).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	} else {
		_, err = db.Exec("INSERT INTO finished_files(hash) values (?)", asciichecksum)
		if err != nil {
			return false
		} else {
			return true
		}
	}
}

func ParseLogFiles() error {
	regex, err := regexp.Compile(shared.MyConfig.LogLineRegex)
	if err != nil {
		return fmt.Errorf("error compiling log line regex: %v", err)
	}

	err = filepath.Walk(shared.MyConfig.LogFilePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), shared.MyConfig.LogFileFilter) {

			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file %s: %v", path, err)
			}
			defer file.Close()

			hash := sha256.New()
			_, err = io.Copy(hash, file)
			if err != nil {
				fmt.Printf("Error reading file content: %v\n", err)
				return err
			}

			checksum := hash.Sum(nil)
			asciichecksum := fmt.Sprintf("%x", checksum)
			proceed := CheckifFileHandled(asciichecksum)
			if _, err := file.Seek(0, 0); err != nil {
				return fmt.Errorf("error resetting file position: %v", err)
			}
			//fmt.Printf("[DEBUG] path %s\n", path)
			if proceed {

				scanner := bufio.NewScanner(file) // Use file directly, as it's already an io.Reader
				if strings.HasSuffix(path, ".gz") {
					gzReader, err := gzip.NewReader(file)
					if err != nil {
						return fmt.Errorf("error creating gzip reader for file %s: %v", path, err)
					}
					defer gzReader.Close()
					scanner = bufio.NewScanner(gzReader) // Update scanner for gzipped files
				}
				//fmt.Printf("debug: hier geraak ik\n%+v\n", reader)
				for scanner.Scan() {
					line := scanner.Text()
					if err != nil {
						if err != io.EOF {
							fmt.Printf("Error reading line: %v\n", err)
							return err
						}
						fmt.Printf("Error reading line: %v\n", err)
						break
					}
					match := regex.FindStringSubmatch(line)
					if match == nil {
						fmt.Printf("no matchy matchy: %s\n", line)
						continue
					}

					var logLine shared.LogLine
					for i, name := range regex.SubexpNames() {
						if i != 0 && name != "" {
							switch name {
							case "RemoteIP":
								logLine.Remote_ip = match[i]
							case "RemoteHostname":
								logLine.Remote_hostame = match[i]
							case "Timestamp":
								t, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[i])
								if err != nil {
									return fmt.Errorf("error parsing timestamp: %v", err)
								}
								logLine.Timestamp = t
							case "RequestMethod":
								logLine.Request_method = match[i]
							case "RequestURI":
								logLine.Request_URI = match[i]
							case "HTTPProtocol":
								logLine.HTTP_Protocol = match[i]
							case "StatusCode":
								fmt.Sscanf(match[i], "%d", &logLine.Status_code)
							case "BytesSent":
								fmt.Sscanf(match[i], "%d", &logLine.Bytes_sent)
							case "RemoteLogname":
								logLine.Remote_logname = match[i]
							case "AuthUser":
								logLine.Remote_user = match[i]
							case "Agent":
								logLine.User_Agent = match[i]
							}
						}
					}
					_ = loglines.AddLogLine(logLine)
				}
				if err := scanner.Err(); err != nil {
					fmt.Printf("Error reading lines: %v\n", err)
					return err
				}
				_ = SetFileHandled(asciichecksum)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through files in %s: %v", shared.MyConfig.LogFilePath, err)
	}

	return nil
}
