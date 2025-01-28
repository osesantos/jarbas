package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	apiURLKey         = "api_url: "
	requestHeadersKey = "request_headers: "
	apiKey            = "api_key: "
	modelKey          = "model: "
	saveMessagesKey   = "save_messages: "
)

const configFile = "/.jarbasrc"

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	path := filepath.Join(homeDir, configFile)

	f, err := os.OpenFile(path, os.O_WRONLY, 0o644)
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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", apiKey, key))
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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", modelKey, model))
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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", saveMessagesKey, save))
	if err != nil {
		return err
	}
	return nil
}

func _writeApiUrl(f *os.File) error {
	fmt.Println("What's the API URL: ")
	// Read input from the user
	apiUrl := ""
	_, err := fmt.Scanln(&apiUrl)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", apiURLKey, apiUrl))
	if err != nil {
		return err
	}
	return nil
}

func _getFile() (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	path := filepath.Join(homeDir, configFile)

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Unable to find config file, please run 'jarbas init' to create one.")
		return nil, err
	}

	return file, nil
}

func GetKey() (string, error) {
	file, err := _getFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, apiKey) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, apiKey)
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetModel() (string, error) {
	file, err := _getFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, modelKey) {
			// Remove the "api-key: " prefix from the line
			apiKey := strings.TrimPrefix(line, modelKey)
			return apiKey, nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetSaveMessages() (bool, error) {
	file, err := _getFile()
	if err != nil {
		return false, err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, saveMessagesKey) {
			option := strings.TrimPrefix(line, saveMessagesKey)

			if option == "n" {
				return false, nil
			}

			if option == "y" {
				return true, nil
			}

			return false, fmt.Errorf("invalid option: %s", option)
		}
	}

	// Return an error if the file is empty
	return false, bufio.ErrFinalToken
}

func GetAPIUrl() (string, error) {
	file, err := _getFile()
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, apiURLKey) {
			return strings.TrimPrefix(line, apiURLKey), nil
		}
	}

	// Return an error if the file is empty
	return "", bufio.ErrFinalToken
}

func GetRequestHeaders() (map[string]string, error) {
	file, err := _getFile()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	headers := map[string]string{}

	// Read the first line of the file
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, requestHeadersKey) {
			headersLine := strings.TrimPrefix(line, requestHeadersKey)
			headersList := strings.Split(headersLine, ",")
			for _, header := range headersList {
				headerParts := strings.Split(header, ":")
				headers[headerParts[0]] = headerParts[1]
			}
			return headers, nil
		}
	}

	// Return an error if the file is empty
	return nil, bufio.ErrFinalToken
}
