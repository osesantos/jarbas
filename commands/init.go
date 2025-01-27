package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const api_key = "api_key: "
const model_key = "model: "
const save_messages_key = "save_messages: "
const config_file = "/.jarbasrc"

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	path := filepath.Join(homeDir, config_file)

	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		f, err = os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer f.Close()

	fmt.Println("Starting config init...")

	err = _writeKey(f)
	if err != nil {
		fmt.Println(err)
	}

	err = _writeModel(f)
	if err != nil {
		fmt.Println(err)
	}

	err = _writeSaveMessages(f)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("File " + f.Name() + " created!")
}

func _writeKey(f *os.File) error {
	fmt.Println("What's the API key: ")
	// Read input from the user
	key := ""
	_, err := fmt.Scanln(&key)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", api_key, key))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func _writeModel(f *os.File) error {
	fmt.Println("What's the model: ")
	// Read input from the user
	model := ""
	_, err := fmt.Scanln(&model)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model_key, model))
	if err != nil {
		return err
	}
	return nil
}

func _writeSaveMessages(f *os.File) error {
	fmt.Println("Do you want to save the messages? (y/n): ")
	// Read input from the user
	save := ""
	_, err := fmt.Scanln(&save)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", save_messages_key, save))
	if err != nil {
		return err
	}
	return nil
}

func GetKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	path := filepath.Join(homeDir, config_file)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return "", err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, api_key) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, api_key)
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetModel() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	path := filepath.Join(homeDir, config_file)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return "", err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, model_key) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, model_key)
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetSaveMessages(string, error) (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	path := filepath.Join(homeDir, config_file)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return false, err
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, save_messages_key) {
			option := strings.TrimPrefix(line, save_messages_key)

			if option == "n" {
				return false, nil
			}

			if option == "y" {
				return true, nil
			}

			return false, fmt.Errorf("Invalid option: %s", option)
		}
	}

	// Return an error if the file is empty
	return false, bufio.ErrFinalToken
}
