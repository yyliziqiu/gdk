package xboot

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type CommandConfig struct {
	Name string
	Desc string
	Long string
	Func func(cmd *cobra.Command, args []string)
}

func ExecuteCommands(app *App, f func(app *App) []*CommandConfig) {
	root := createRootCommand(app)

	for _, cfg := range f(app) {
		if cfg.Long == "" {
			cfg.Long = cfg.Desc
		}
		root.AddCommand(&cobra.Command{
			Use:   cfg.Name,
			Short: cfg.Desc,
			Long:  cfg.Long,
			Run:   cfg.Func,
		})
	}

	if err := root.Execute(); err != nil {
		fmt.Printf("Execute command failed, error: %v.", err)
		os.Exit(1)
	}
}

func createRootCommand(app *App) *cobra.Command {
	var config string

	root := &cobra.Command{
		Version: app.Version,
		Use:     app.Command,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Use %s.bin -h or --help for help.\n", app.Command)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// 覆盖配置文件路径
			if strings.TrimSpace(config) != "" {
				app.ConfigFile = config
			}
			// 初始化配置
			if err := app.InitConfig(); err != nil {
				fmt.Printf("Init config failed, error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	root.PersistentFlags().StringVarP(&config, "config", "c", "", "config file path")

	return root
}

func ExecuteCommand(app *App, commands ...func(app *App) *cobra.Command) {
	root := createRootCommand(app)

	for _, cmd := range commands {
		root.AddCommand(cmd(app))
	}

	if err := root.Execute(); err != nil {
		fmt.Printf("Execute command failed, error: %v.", err)
		os.Exit(1)
	}
}
