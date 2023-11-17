package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mheers/inline-age/crypt"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
)

var (
	identityFile string

	decryptCmd = &cobra.Command{
		Use:     "decrypt [chiffre]",
		Short:   "decrypts a string with an identity file or a password",
		Long:    ``,
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			chiffre := args[0]

			if chiffre == "-" {
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					chiffre = scanner.Text()
				} else {
					return fmt.Errorf("no input")
				}
			}

			var plaintext string
			var err error
			if password != "" {
				plaintext, err = crypt.DecryptFromPassword(chiffre, password)
			} else {
				plaintext, err = crypt.DecryptFromIdentityFile(chiffre, identityFile)
			}
			if err != nil {
				return err
			}
			fmt.Println(plaintext)

			return nil
		},
	}
)

func init() {
	decryptCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", helpers.PrivateKeyPath(), "")
	decryptCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "")
}
