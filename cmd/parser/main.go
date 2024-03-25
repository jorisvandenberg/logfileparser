package main

import (
	"fmt"
	"logfileparser/shared/config"
	"logfileparser/shared/db/initialise"
	"logfileparser/shared/parser"
)

func main() {
	config.LoadMyConfig()
	err := initialise.InitialiseDb()
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		return
	}
	_ = parser.ParseLogFiles()
}
