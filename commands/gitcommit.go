package commands

import (
	"jarbas-go/main/actions"
	"jarbas-go/main/prompts"
	"jarbas-go/main/settings"
	"os/exec"
)

func GetGitCommit(settings settings.Settings) (string, error) {
	output, err := getdiff()
	if err != nil {
		return "", err
	}
	if len(output) == 0 {
		return "No changes found...", nil
	}

	prompt := prompts.GetGitCommit(string(output))
	response := actions.SingleQuestion(prompt, settings, "")

	if response == "" {
		return "", nil
	}
	return response, nil
}

func getdiff() ([]byte, error) {
	output, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
