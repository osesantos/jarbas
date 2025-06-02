package utils

import (
	"fmt"
	"jarbas-go/main/model"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

func extractTimestamp(filename string) string {
	// File format: <uuid>-<timestamp>.json
	// Extract the timestamp part
	block := strings.Replace(filename, ".json", "", 1)
	parts := strings.Split(block, "-")

	if len(parts) < 2 {
		return ""
	}

	// Lets check if there is something after
	smallerParts := strings.Split(parts[len(parts)-1], " ")
	if len(smallerParts) > 1 {
		return smallerParts[0] // Return the first part as the timestamp
	}

	return parts[len(parts)-1] // Return the last part as the timestamp
}

func OrderFilesByTime(files []string) []string {
	// File format: <uuid>-<timestamp>.json
	// sort by timestamp
	filesCopy := make([]string, len(files))
	copy(filesCopy, files)

	sort.Slice(filesCopy, func(i, j int) bool {
		ti := extractTimestamp(filesCopy[i])
		tj := extractTimestamp(filesCopy[j])

		return ti > tj // Sort in descending order
	})

	return filesCopy
}

func AddDateTimeToFiles(files []string) []string {
	for i, file := range files {
		// Extract the timestamp part from the filename
		timestamp := extractTimestamp(file)

		if timestamp == "" {
			continue // Skip files without a valid timestamp
		}

		unixSec, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			continue // Skip files with invalid timestamps
		}

		t := time.Unix(unixSec, 0)

		files[i] = fmt.Sprintf("%s (%s)", file, t.Format("02-01-2006 15:04:05"))
	}
	return files
}

func AddTitleToFiles(files []string, parser func(file string) (model.Chat, error)) []string {
	for i, file := range files {
		conversation, err := parser(file)
		if err != nil {
			// Since we cant find the title for this file, we will just skip it
			continue // Skip files with invalid JSON
		}

		files[i] = fmt.Sprintf("%s %s", file, conversation.Title)
	}
	return files
}

func CleanFileName(fileName string) string {
	// Remove the "(date time)" part from the filename
	fileName = strings.Split(fileName, " ")[0]

	return fileName
}
