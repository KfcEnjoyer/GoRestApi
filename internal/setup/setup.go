package setup

import (
	"fmt"
	"log"
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
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				log.Println(err)
			}
			fmt.Println("Path created")
		}
	}

	for filePath, contend := range defaultFiles {
		fullPath := filepath.Join(basePath, filePath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if err := os.WriteFile(fullPath, []byte(contend), 0644); err != nil {
				log.Println(err)
			}
		}
	}
}
