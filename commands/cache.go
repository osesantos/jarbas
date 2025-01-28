package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const cacheDir = "/.local/share/jarbas"

func GetCacheDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	path := filepath.Join(homeDir, cacheDir)
	return path
}

func CreateCacheDir() error {
	path := GetCacheDir()
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func SaveConversation(messages []map[string]interface{}) error {
	path := GetCacheDir()
	// get random uuid for the conversation file
	id := uuid.NewString()
	id = strings.ReplaceAll(id, "-", "")

	// get timestamp for the conversation file
	timestamp := time.Now().Unix()

	fileName := fmt.Sprintf("%s-%d.json", id, timestamp)

	// check if dir exists if not create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := CreateCacheDir()
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	// convert to json string
	jsonData, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("%v", string(jsonData)))
	if err != nil {
		return err
	}

	return nil
}

func GetConversations() ([]map[string]interface{}, error) {
	path := GetCacheDir()
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var conversations []map[string]interface{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(path, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		var conversation []map[string]interface{}
		_, err = fmt.Fscan(file, &conversation)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conversation...)
	}

	return conversations, nil
}
