package storage

import (
	"GoRestApi/internal/api"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var SaveFilePath = "save/save_req.json"
var savedReq = make(map[string][]api.Req)

func EnsureFile() {
	saveDir := filepath.Dir(SaveFilePath)
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.MkdirAll(saveDir, os.ModePerm)
		fmt.Println("Path created")
	}

	if _, err := os.Stat(SaveFilePath); os.IsNotExist(err) {
		os.WriteFile(SaveFilePath, []byte("{}"), 0644)
	}
}

func SaveRequest(name string, req api.Req) error {
	EnsureFile()

	data, err := LoadRequests()
	if err != nil {
		return err
	}

	data[name] = append(data[name], req)

	newData, err := json.MarshalIndent(&data, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(SaveFilePath, newData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("request saved ")
	return nil
}

func DeleteRequest(name string, index int) error {
	EnsureFile()

	data, err := LoadRequests()
	if err != nil {
		return err
	}

	requests, exists := data[name]
	if !exists {
		return errors.New("request name does not exist")
	}

	if index < 0 || index >= len(requests) {
		return errors.New("invalid request index")
	}

	if len(requests) == 1 {
		delete(data, name)
	} else {
		data[name] = append(requests[:index], requests[index+1:]...)
	}

	newData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(SaveFilePath, newData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("request deleted")
	return nil
}

func LoadRequests() (map[string][]api.Req, error) {
	EnsureFile()

	read, err := os.ReadFile(SaveFilePath)
	if err != nil {
		return nil, err
	}

	var data map[string][]api.Req
	err = json.Unmarshal(read, &data)
	fmt.Println(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
