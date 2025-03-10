package main

import (
	"fmt"
	"log"
	"os"

	"jarbas-go/main/actions"
	"jarbas-go/main/agents"
	"jarbas-go/main/commands"
	"jarbas-go/main/model"
	"jarbas-go/main/vendors"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Jarbas",
		Usage:   "A chatGPT cli implementation that uses API to have the ChatGPT to get help on the terminal",
		Version: "0.1.0",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Generate ~/.jarbasrc file.",
				Action: func(cCtx *cli.Context) error {
					commands.Init()
					return nil
				},
			},
			{
				Name:    "chat",
				Aliases: []string{"c"},
				Usage:   "Start a chat with jarbas",
				Action: func(cCtx *cli.Context) error {
					settings, err := model.GetSettings()
					if err != nil {
						return err
					}

					err = commands.Chat(settings, nil)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "continue-chat",
				Aliases: []string{"cc"},
				Usage:   "Continue a chat with jarbas",
				Action: func(cCtx *cli.Context) error {
					settings, err := model.GetSettings()
					if err != nil {
						return err
					}

					err = commands.ContinueChat(settings)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "agent",
				Aliases: []string{"a"},
				Usage:   "Run an agent",
				Action: func(cCtx *cli.Context) error {
					settings, err := model.GetSettings()
					if err != nil {
						return err
					}

					err = agents.RunAgent("", settings)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "~/.jarbasrc",
				Aliases: []string{"i"},
				Usage:   "Input configuration from `FILE`",
			},
		},
		Action: func(cCtx *cli.Context) error {
			settings, err := model.GetSettings()
			if err != nil {
				return err
			}

			question := cCtx.Args().Get(1)
			response, err := actions.SingleQuestion(question, settings, vendors.SoftwareEngineer())
			if err != nil {
				return err
			}
			fmt.Println(commands.QuestionPrompt + question)
			fmt.Println(commands.DefaultPrompt + response)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
