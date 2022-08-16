package um

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

//go:embed migrations/*.sql
var fs embed.FS

func CreateUmMigrationsFile(dir string) error {
	const folder = "migrations"
	entries, err := fs.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		content, err := fs.ReadFile(fmt.Sprintf("%s/%s", folder, entry.Name()))
		if err != nil {
			return err
		}

		target := fmt.Sprintf("%s/%s", dir, entry.Name())
		if err := os.WriteFile(target, content, 0666); err != nil {
			panic(err)
		}
	}
	return nil
}

func readQueries() ([]string, []string, error) {
	const folder = "migrations"
	fname := "001_create_user_roles_tbl.sql"
	data, err := fs.ReadFile(fmt.Sprintf("%s/%s", folder, fname))
	if err != nil {
		return nil, nil, err
	}
	contents := string(data)
	const sep = "---- create above / drop below ----"
	parts := strings.Split(contents, sep)
	upAll := strings.Split(strings.TrimSpace(parts[0]), "--#")
	downAll := strings.Split(strings.TrimSpace(parts[1]), "--#")
	var up, down []string
	for _, q := range upAll {
		up = append(up, strings.TrimSpace(q))
	}
	for _, q := range downAll {
		down = append(down, strings.TrimSpace(q))
	}
	return up, down, nil
}
