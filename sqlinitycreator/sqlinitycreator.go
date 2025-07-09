package sqlinitycreator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mwac-dev/sqlinity/sqlinitytypes"
)

func CreateMigrationFile(config sqlinitytypes.Config, rawName string) error {
	files, err := os.ReadDir(config.SqlFolder)
	if err != nil {
		return err
	}
	maxID := 0
	re := regexp.MustCompile(`^(\d+)_.*\.up\.sql$`)
	for _, file := range files {
		match := re.FindStringSubmatch(file.Name())
		if len(match) == 2 {
			num, err := strconv.Atoi(match[1])
			if err == nil && num > maxID {
				maxID = num
			}
		}
	}

	newID := fmt.Sprintf("%03d", maxID+1)
	slug := strings.ToLower(strings.ReplaceAll(rawName, " ", "_"))

	upFile := filepath.Join(config.SqlFolder, fmt.Sprintf("%s_%s.up.sql", newID, slug))
	downFile := filepath.Join(config.SqlFolder, fmt.Sprintf("%s_%s.down.sql", newID, slug))

	err = os.WriteFile(upFile, []byte("-- TODO: Write your up migration SQL here\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	err = os.WriteFile(downFile, []byte("-- TODO: Write your down migration SQL here\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Created migration files:\n- %s\n- %s\n", upFile, downFile)
	return nil

}
