package summarizer

import (
	"fmt"

	"jarbas-go/main/actions"
	"jarbas-go/main/commands"
	"jarbas-go/main/prompts"
	"jarbas-go/main/settings"
	"jarbas-go/main/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/osesantos/resulto"
)

type Options struct {
	URL string `json:"url"`
}

func _prompt(scrapedText string) string {
	return fmt.Sprintf(
		`
		SYSTEM PROMPT:
		YOU are a PROFESSIONAL WRITER, and I need you to summarize the following text
		in a few sentences. Try to keep the main points and the most important details.
		--------------------------------
		USER PROMPT:
		%s`, scrapedText)
}

func _getURL() resulto.Result[string] {
	question := ""
	prompt := &survey.Input{
		Message: "url to summarize: ",
	}
	err := survey.AskOne(prompt, &question)
	if err != nil {
		return resulto.Failure[string](err)
	}

	return resulto.Success(question)
}

func GetOptions() Options {
	url := _getURL().Unwrap()

	return Options{
		URL: url,
	}
}

func Run(options Options, settings settings.Settings) {
	scarpedText := utils.ScrapeText(options.URL).Unwrap()

	prompt := _prompt(scarpedText)
	response := actions.SingleQuestion(prompt, settings, prompts.ProfessionalWriter())

	fmt.Println(commands.DefaultPrompt + response)
}
