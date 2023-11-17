package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mheers/inline-age/file"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
)

var (
	encryptFileCmd = &cobra.Command{
		Use:     "encrypt-file [file] [path] [plaintext]",
		Short:   "encrypts a string ath path in a file with PublicSecret stored in that file - decrypted with identity file and writes the reencrypted string back",
		Long:    ``,
		Aliases: []string{"ef"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) < 3 {
				return fmt.Errorf("not enough arguments")
			}

			filename := args[0]
			path := args[1]
			plaintext := args[2]

			if plaintext == "-" {
				// read from stdin
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					plaintext = scanner.Text()
				} else {
					return fmt.Errorf("no input")
				}
			}

			err := file.EncryptJSONFilePath(filename, path, plaintext, identityFile)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	encryptFileCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", helpers.PrivateKeyPath(), "")
}
