package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run create_migration.go <migration_name>")
		return
	}

	migrationName := os.Args[1]
	migrationName = strings.ToLower(strings.ReplaceAll(migrationName, " ", "_"))

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")
	baseFileName := fmt.Sprintf("%s_%s", timestamp, migrationName)

	// Create migrations folder if not exists
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		err := os.Mkdir("migrations", os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating migrations folder: %v\n", err)
			return
		}
	}

	// Create .up.sql file
	upFilePath := fmt.Sprintf("migrations/%s.up.sql", baseFileName)
	upFile, err := os.Create(upFilePath)
	if err != nil {
		fmt.Printf("Error creating up migration file: %v\n", err)
		return
	}
	defer upFile.Close()

	// Create .down.sql file
	downFilePath := fmt.Sprintf("migrations/%s.down.sql", baseFileName)
	downFile, err := os.Create(downFilePath)
	if err != nil {
		fmt.Printf("Error creating down migration file: %v\n", err)
		return
	}
	defer downFile.Close()

	fmt.Printf("✅ Migration files created:\n- %s\n- %s\n", upFilePath, downFilePath)
}
