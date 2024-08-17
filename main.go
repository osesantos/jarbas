package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"jarbas-go/main/actions"
	"jarbas-go/main/commands"
	"log"
	"os"
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
					model, err := commands.GetModel()
					if err != nil {
						return err
					}

					err = commands.Chat(key, model)
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
