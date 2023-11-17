package cmd

import (
	"fmt"

	"github.com/mheers/inline-age/file"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
)

var (
	paths []string

	reencryptFileCmd = &cobra.Command{
		Use:     "reencrypt-file [file]",
		Short:   "reencrypts a a file",
		Long:    ``,
		Aliases: []string{"ref"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			fileName := args[0]

			if recipientFile == "" {
				cmd.Println("no recipient file set, using recipients from secret file")
				recipientFile = fileName
			}

			if len(paths) > 0 {
				err := file.ReEncryptJSONFilePaths(fileName, paths, identityFile, recipientFile)
				if err != nil {
					return err
				}
			} else {
				err := file.ReencryptJSONFile(fileName, identityFile, recipientFile)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
)

func init() {
	reencryptFileCmd.PersistentFlags().StringVarP(&recipientFile, "recipient-file", "r", recipientFileDefault, "if not set the recipients defined in the secret file will be used")
	reencryptFileCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", helpers.PrivateKeyPath(), "")
	reencryptFileCmd.PersistentFlags().StringArrayVarP(&paths, "paths", "p", []string{}, "")
}
