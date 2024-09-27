package cmd

import (
	"fmt"

	"github.com/mheers/inline-age/file"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
)

var (
	initFileCmd = &cobra.Command{
		Use:     "init-file [file]",
		Short:   "inits a file with PublicSecret stored in that file",
		Long:    ``,
		Aliases: []string{"if"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			filename := args[0]

			// TODO: check if already initialized - then abort

			err := file.InitJSONFile(filename, recipientFile)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	initFileCmd.PersistentFlags().StringVarP(&recipientFile, "recipient-file", "r", "", "")
}
