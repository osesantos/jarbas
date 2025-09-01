package main

import (
	"fmt"
	"log"
	"os"

	"jarbas-go/main/actions"
	"jarbas-go/main/agents"
	"jarbas-go/main/commands"
	"jarbas-go/main/prompts"
	"jarbas-go/main/settings"

	"github.com/charmbracelet/glamour"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Jarbas",
		Usage:   "A chatGPT cli implementation that uses API to have the ChatGPT to get help on the terminal",
		Version: "1.2.0",
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
					settings := settings.GetSettings(cCtx)
					err := commands.Chat(settings, nil, false)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "mcp",
				Aliases: []string{"mcp"},
				Usage:   "Do a single question using MCP backend",
				Action: func(cCtx *cli.Context) error {
					settings := settings.GetSettings(cCtx)

					renderer, err := glamour.NewTermRenderer(
						glamour.WithStandardStyle("dark"),
						glamour.WithEmoji(),
						glamour.WithWordWrap(0),
						glamour.WithTableWrap(true),
					)
					if err != nil {
						return fmt.Errorf("error initializing glamour renderer: %w", err)
					}

					query := cCtx.Args().First()
					spinner, _ := pterm.DefaultSpinner.Start("Asking MCP...")
					answer := commands.McpQuery(settings, query)
					spinner.Success("MCP answered:")
					out, err := renderer.Render(answer)
					if err != nil {
						fmt.Println("Error rendering message:", err)
					} else {
						fmt.Println(out)
					}
					return nil
				},
			},
			{
				Name:    "continue-chat",
				Aliases: []string{"cc"},
				Usage:   "Continue a chat with jarbas",
				Action: func(cCtx *cli.Context) error {
					settings := settings.GetSettings(cCtx)
					err := commands.ContinueChat(settings)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "oneliner",
				Aliases: []string{"ol"},
				Usage:   "Ask jarbas to generate a oneliner for you, always start with 'linux:' or 'windows:' to specify the OS",
				Action: func(cCtx *cli.Context) error {
					settings := settings.GetSettings(cCtx)
					oneliner, err := commands.GetOneLiner(settings, cCtx.Args().First())
					if err != nil {
						return err
					}
					fmt.Println(commands.DefaultPrompt + oneliner)
					return nil
				},
			},
			{
				Name:    "gitcommit",
				Aliases: []string{"gc"},
				Usage:   "Ask jarbas to generate a git commit command for you based on the current dir cached changes",
				Action: func(cCtx *cli.Context) error {
					settings := settings.GetSettings(cCtx)
					commitCommand, err := commands.GetGitCommit(settings)
					if err != nil {
						return err
					}
					fmt.Println(commands.DefaultPrompt + commitCommand)
					return nil
				},
			},
			{
				Name:    "agent",
				Aliases: []string{"a"},
				Usage:   "Run an agent",
				Action: func(cCtx *cli.Context) error {
					settings := settings.GetSettings(cCtx)
					agents.RunAgent("", settings)
					return nil
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "sync saved chats to the gomind server",
				Action: func(cCtx *cli.Context) error {
					commands.SyncChats()
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
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "Specify the model to use, e.g., gpt-3.5-turbo",
			},
		},
		Action: func(cCtx *cli.Context) error {
			settings := settings.GetSettings(cCtx)
			question := cCtx.Args().Get(1)
			response := actions.SingleQuestion(question, settings, prompts.SoftwareEngineer())
			fmt.Println(commands.QuestionPrompt + question)
			fmt.Println(commands.DefaultPrompt + response)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
