package sqlinityparser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mwac-dev/sqlinity/sqlinitytypes"
)

func ParseMigrations(config sqlinitytypes.Config) ([]sqlinitytypes.Migration, error) {
	var migrations []sqlinitytypes.Migration

	files, err := os.ReadDir(config.SqlFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		baseName := strings.TrimSuffix(name, ".up.sql") // Remove the ".up.sql" suffix and just keep the base name
		parts := strings.SplitN(baseName, "_", 2)       // Split the base name into ID and Name
		if len(parts) != 2 {
			fmt.Printf("Skipping file %s: does not follow naming convention (ID_Name.up.sql)\n", name)
			continue
		}
		id := parts[0]
		migrationName := parts[1]

		upPath := filepath.Join(config.SqlFolder, name)
		downPath := strings.Replace(upPath, ".up.sql", ".down.sql", 1)

		upSQL, err := os.ReadFile(upPath)
		if err != nil {
			return nil, fmt.Errorf("error reading up SQL file %s: %w", upPath, err)
		}

		// Optional down sql file
		downSQL := []byte{}
		if _, err := os.Stat(downPath); err == nil {
			downSQL, err = os.ReadFile(downPath)
			if err != nil {
				return nil, fmt.Errorf("error reading down SQL file %s: %w", downPath, err)
			}
		}

		m := sqlinitytypes.Migration{
			ID:      id,
			Name:    migrationName,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}
