package cmd

import (
	"fmt"

	"github.com/mheers/inline-age/file"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
)

var (
	decryptFileCmd = &cobra.Command{
		Use:     "decrypt-file [file] [path]",
		Short:   "decrypts a string at path in a file with PublicSecret stored in that file - decrypted with identity file and prints it",
		Long:    ``,
		Aliases: []string{"df"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) < 2 {
				return fmt.Errorf("not enough arguments")
			}

			filename := args[0]
			path := args[1]

			plaintext, err := file.DecryptJSONFilePath(filename, path, identityFile)
			if err != nil {
				return err
			}

			fmt.Println(plaintext)

			return nil
		},
	}
)

func init() {
	decryptFileCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", helpers.PrivateKeyPath(), "")
}
