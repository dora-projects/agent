package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ween/internal"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		serverRun()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int("port", 6000, "server port")

	err := viper.BindPFlags(serverCmd.Flags())
	if err != nil {
		panic(err)
	}
}

func serverRun() {
	internal.Proxy()
}
