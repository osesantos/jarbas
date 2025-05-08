package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/osesantos/resulto"
)

const ConfigFile = "/.jarbasrc"

func GetSettingsFile() resulto.Result[*os.File] {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return resulto.Failure[*os.File](err)
	}

	path := filepath.Join(homeDir, ConfigFile)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return resulto.Failure[*os.File](err)
	}

	return resulto.Success(file)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func FileNotEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.Size() > 0
}

func OrderFilesByTime(files []string) []string {
	// File format: <uuid>-<timestamp>.json
	// sort by timestamp
	sort.Slice(files, func(i, j int) bool {
		timeString := strings.Split(files[i], "-")[1]
		timeString = strings.Replace(timeString, ".json", "", 1)

		fmt.Println("Parsing: " + timeString)
		iTime, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			fmt.Println(err)
			return false
		}

		timeString = strings.Split(files[j], "-")[1]
		timeString = strings.Replace(timeString, ".json", "", 1)

		fmt.Println("Parsing: " + timeString)
		jTime, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			fmt.Println(err)
			return false
		}

		return iTime.Before(jTime)
	})

	return files
}
