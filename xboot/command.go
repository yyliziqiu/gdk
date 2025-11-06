package xboot

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var _rootCommand *cobra.Command

func ExecuteCommand(app *App, commands ...func(app *App) *cobra.Command) {
	initRootCommand(app)

	for _, command := range commands {
		_rootCommand.AddCommand(command(app))
	}

	err := _rootCommand.Execute()
	if err != nil {
		fmt.Printf("Execute command failed, error: %v.", err)
		os.Exit(1)
	}
}

func initRootCommand(app *App) {
	var config string

	_rootCommand = &cobra.Command{
		Version: app.Version,
		Use:     app.Name,
		Run: func(command *cobra.Command, args []string) {
			fmt.Printf("Use %s.bin -h or --help for help.\n", app.Name)
		},
		PersistentPreRun: func(command *cobra.Command, args []string) {
			if strings.TrimSpace(config) != "" {
				app.ConfigPath = config
			}
			err := app.InitConfig()
			if err != nil {
				fmt.Printf("Init config failed, error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	_rootCommand.PersistentFlags().StringVarP(&config, "config", "c", "", "config file path")
}
