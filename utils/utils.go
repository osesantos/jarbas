package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFile = "/.jarbasrc"

func GetSettingsFile() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	path := filepath.Join(homeDir, ConfigFile)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return nil, err
	}

	return file, nil
}
