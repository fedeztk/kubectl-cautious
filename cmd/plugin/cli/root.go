package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/fedeztk/kubectl-cautious/pkg/config"
	"github.com/fedeztk/kubectl-cautious/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd returns the root command for the plugin
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cautious [kubectl command]`,
		Short: "A kubectl plugin to prevent accidental changes",
		Long: `A kubectl plugin to prevent accidental changes. See the automatically generated configuration file
at $HOME/.kube/cautious.yaml after running the plugin for the first time. Supports regexes.`,
		SilenceErrors:      true,
		SilenceUsage:       true,
		DisableFlagParsing: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// INFO: if this is enabled, the help command of kubectl will not work
			// it is best to use the kubectl krew info cautious command
			// if len(args) == 1 {
			// 	if args[0] == "--help" || args[0] == "-h" {
			// 		cmd.Help()
			// 		versionInfo := version.GetVersion()
			// 		fmt.Printf("%s", versionInfo.ToString())
			// 		return nil
			// 	}
			// }

			conf, err := config.ReadConfig()
			if err != nil {
				return err
			}

			if err := plugin.RunPlugin(conf, args); err != nil {
				return err
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

// InitAndExecute initializes the plugin and executes the root
func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		// log only if the error is not a kubectl error
		if errors.Is(err, plugin.ErrKubectl) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}

func initConfig() {
	log.SetReportTimestamp(false)
	viper.AutomaticEnv()
}
