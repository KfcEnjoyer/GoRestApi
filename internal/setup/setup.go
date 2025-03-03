package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

var basePath = "./"

var directories = []string{"save", "logs"}

var defaultFiles = map[string]string{
	"save/save_req.json": "{}",
	"logs/app.log":       "[]",
}

func EnsureDirectoriesAndFiles() {
	for _, dir := range directories {
		path := filepath.Join(basePath, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, os.ModePerm)
			fmt.Println("Path created")
		}
	}

	for filePath, contend := range defaultFiles {
		fullPath := filepath.Join(basePath, filePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			os.WriteFile(fullPath, []byte(contend), 0644)
		}
	}
}
