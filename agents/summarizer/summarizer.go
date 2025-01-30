package summarizer

import (
	"fmt"
	"jarbas-go/main/actions"
	"jarbas-go/main/commands"
	"jarbas-go/main/model"
	"jarbas-go/main/utils"

	"github.com/AlecAivazis/survey/v2"
)

type Options struct {
	URL string `json:"url"`
}

func _prompt(scrapedText string) string {
	return fmt.Sprintf(
		`YOU are a PROFESSIONAL WRITER, and I need you to summarize the following text
		in a few sentences. Try to keep the main points and the most important details.
		--------------------------------
		%s`, scrapedText)
}

func _getUrl() (string, error) {
	question := ""
	prompt := &survey.Input{
		Message: "url to summarize: ",
	}
	err := survey.AskOne(prompt, &question)
	if err != nil {
		return "", err
	}

	return question, nil
}

func GetOptions() (Options, error) {
	url, err := _getUrl()
	if err != nil {
		return Options{}, err
	}

	return Options{
		URL: url,
	}, nil
}

func Run(options Options, settings model.Settings) error {
	scarpedText, err := utils.ScrapeText(options.URL)
	if err != nil {
		return err
	}

	prompt := _prompt(scarpedText)

	response, err := actions.SingleQuestion(prompt, settings)
	if err != nil {
		return err
	}

	fmt.Println(commands.DefaultPrompt + response)

	return nil
}
