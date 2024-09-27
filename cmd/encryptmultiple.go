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
	encryptMultipleCmd = &cobra.Command{
		Use:     "encrypt-multiple [plaintexts...]",
		Short:   "encrypts multiple strings with a recipient file and prints the public secret",
		Long:    ``,
		Aliases: []string{"em"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			plaintexts := []string{}

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			if args[0] == "-" {
				// read from stdin
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					plaintexts = append(plaintexts, scanner.Text())
				}

				if len(plaintexts) == 0 {
					return fmt.Errorf("not enough arguments")
				}
			} else {
				plaintexts = args[0:]
			}

			chiffres, publicSecret, err := crypt.EncryptMultipleCommon(plaintexts, recipientFile)
			if err != nil {
				return err
			}

			type iaConfig struct {
				PublicSecret string
			}

			type result struct {
				IAConfig iaConfig `json:"__ia_config__"`
				Chiffres []string
			}

			r := result{
				IAConfig: iaConfig{
					PublicSecret: publicSecret,
				},
				Chiffres: chiffres,
			}

			return helpers.PrintFormat(r, outputFormatFlag)
		},
	}
)

func init() {
	encryptMultipleCmd.PersistentFlags().StringVarP(&recipientFile, "recipient-file", "r", "", "")
}
