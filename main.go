package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mwac-dev/sqlinity/sqlinitycreator"
	"github.com/mwac-dev/sqlinity/sqlinitygenerator"
	"github.com/mwac-dev/sqlinity/sqlinityparser"
	"github.com/mwac-dev/sqlinity/sqlinitytypes"
)

func loadConfig() (sqlinitytypes.Config, error) {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		return sqlinitytypes.Config{}, err
	}

	var config sqlinitytypes.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return sqlinitytypes.Config{}, err
	}
	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	
	if len(os.Args) >= 3 && os.Args[1] == "create" {
		rawName := strings.Join(os.Args[2:], " ")
		err = sqlinitycreator.CreateMigrationFile(config, rawName)
		if err != nil {
			fmt.Printf("Error creating migration file: %v\n", err)
			return
		}
		return
	}

	fmt.Println("Configuration loaded successfully:")
	fmt.Println("SQL Folder:", config.SqlFolder)
	fmt.Println("Output Folder:", config.OutputFolder)
	fmt.Println("Namespace:", config.Namespace)

	migrations, err := sqlinityparser.ParseMigrations(config)
	if err != nil {
		fmt.Printf("Error parsing migrations: %v\n", err)
		return
	}

	fmt.Println("Migrations parsed successfully: ", len(migrations), "found")
	for _, migration := range migrations {
		fmt.Printf("ID: %s, Name: %s\n", migration.ID, migration.Name)
	}

	err = sqlinitygenerator.GenerateMigrations(config, migrations)
	if err != nil {
		fmt.Printf("Error generating migrations: %v\n", err)
		return
	}
	
	err = sqlinitygenerator.GenerateMigrationRegistry(config, migrations)
	if err != nil {
		fmt.Printf("Error generating migration registry: %v\n", err)
		return
	}
}
