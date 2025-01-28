package commands

import (
	"fmt"
	"jarbas-go/main/model"
	"jarbas-go/main/utils"
	"os"
	"path/filepath"
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	path := filepath.Join(homeDir, utils.ConfigFile)

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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model.GetJsonKey(model.ApiKey), key))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func _writeModel(f *os.File) error {
	fmt.Println("What's the model: ")
	// Read input from the user
	modelValue := ""
	_, err := fmt.Scanln(&modelValue)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model.GetJsonKey(model.Model), modelValue))
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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model.GetJsonKey(model.SaveMessages), save))
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

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model.GetJsonKey(model.ApiURL), apiUrl))
	if err != nil {
		return err
	}
	return nil
}

func _writeVendor(f *os.File) error {
	fmt.Println("What's the vendor: ")
	// Read input from the user
	vendor := ""
	_, err := fmt.Scanln(&vendor)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s%s\n", model.GetJsonKey(model.Vendor), vendor))
	if err != nil {
		return err
	}
	return nil
}
