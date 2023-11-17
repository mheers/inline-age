package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// logLevelFlag describes the verbosity of logs
	logLevelFlag string

	// outputFormatFlag can be json, yaml or table
	outputFormatFlag string

	rootCmd = &cobra.Command{
		Use:   "ia",
		Short: "ia is a command line tool for encrypting and decrypting secrets in files",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevelFlag, "log-level", "l", "error", "possible values are debug, error, fatal, panic, info, trace")
	rootCmd.PersistentFlags().StringVarP(&outputFormatFlag, "output-format", "O", "json", "format [json|yaml|csv]")
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(encryptFileCmd)
	rootCmd.AddCommand(encryptMultipleCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(decryptFileCmd)
	rootCmd.AddCommand(reencryptFileCmd)
	rootCmd.AddCommand(resolveReferencesFileCmd)
	rootCmd.AddCommand(initFileCmd)
	rootCmd.AddCommand(versionCmd)
}
