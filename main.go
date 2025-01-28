package main

import (
	"fmt"
	"log"
	"os"

	"jarbas-go/main/actions"
	"jarbas-go/main/commands"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Jarbas",
		Usage: "A chatGPT cli implementation that uses API to have the ChatGPT to get help on the terminal",
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
					key, err := commands.GetKey()
					if err != nil {
						return err
					}

					model, err := commands.GetModel()
					if err != nil {
						return err
					}

					saveMessages, err := commands.GetSaveMessages()
					if err != nil {
						return err
					}

					err = commands.Chat(key, model, saveMessages, nil)
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
					key, err := commands.GetKey()
					if err != nil {
						return err
					}

					model, err := commands.GetModel()
					if err != nil {
						return err
					}

					saveMessages, err := commands.GetSaveMessages()
					if err != nil {
						return err
					}

					err = commands.ContinueChat(key, model, saveMessages)
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
			key, err := commands.GetKey()
			if err != nil {
				return err
			}

			model, err := commands.GetModel()
			if err != nil {
				return err
			}

			question := cCtx.Args().Get(0)
			response, err := actions.Question(question, key, model)
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
