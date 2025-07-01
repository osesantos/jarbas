package commands

import (
	"encoding/json"
	"fmt"
	"jarbas-go/main/model"
	"os"
	"path/filepath"

	"github.com/osesantos/resulto"
)

const memoryDir = "/.local/share/jarbas/memory"

func GetMemoryDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	path := filepath.Join(homeDir, memoryDir)
	return path
}

func CreateMemoryDir() error {
	path := GetMemoryDir()
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println("Error checking directory:", err)
		return false
	}
	return info.IsDir()
}

func GetMemories() resulto.Result[[]model.Memory] {
	path := GetMemoryDir()
	if !dirExists(path) {
		err := CreateMemoryDir()
		if err != nil {
			return resulto.Failure[[]model.Memory](err)
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return resulto.Failure[[]model.Memory](err)
	}

	var memories []model.Memory
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(path, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return resulto.Failure[[]model.Memory](err)
		}

		defer file.Close()

		var memory model.Memory
		err = json.NewDecoder(file).Decode(&memory)
		if err != nil {
			return resulto.Failure[[]model.Memory](err)
		}

		memories = append(memories, memory)
	}

	return resulto.Success(memories)
}
