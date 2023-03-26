package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	path := filepath.Join(homeDir, "/.jarbasrc")

	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer f.Close()

	fmt.Println("Starting config init...")
	fmt.Println("What's the API key: ")
	// Read input from the user
	key := ""
	fmt.Scanln(&key)

	_, err = f.WriteString("api-key: " + key + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("File " + f.Name() + " created!")
}

func GetKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	path := filepath.Join(homeDir, "/.jarbasrc")

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	if scanner.Scan() {
		line := scanner.Text()

		// Remove the "api-key: " prefix from the line
		apiKey := strings.TrimPrefix(line, "api-key: ")
		return apiKey, nil
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}
