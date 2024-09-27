package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mheers/inline-age/crypt"
	"github.com/mheers/inline-age/helpers"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	recipients    []string
	recipientFile string
	password      string

	encryptCmd = &cobra.Command{
		Use:     "encrypt [plaintext | - for stdin]",
		Short:   "encrypts a string with a recipient file or a list of recipients and prints the encrypted string",
		Long:    ``,
		Aliases: []string{"e"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			plaintext := args[0]

			if plaintext == "-" {
				// read from stdin
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					plaintext = scanner.Text()
				} else {
					return fmt.Errorf("no input")
				}
			}

			var chiffre string
			var err error

			if recipientFile != "" {
				chiffre, err = crypt.EncryptStringWithRecipientFile(plaintext, recipientFile)
			} else if len(recipients) > 0 {
				chiffre, err = crypt.Encrypt(plaintext, recipients)
			} else {
				if password == "" {
					fmt.Print("Enter password: ")
					tty, err := os.Open("/dev/tty")
					if err != nil {
						return fmt.Errorf("error opening /dev/tty: %v", err)
					}
					defer tty.Close()

					fd := int(tty.Fd())
					inputPassword, err := term.ReadPassword(fd)
					if err != nil {
						return fmt.Errorf("error reading password: %v", err)
					}
					fmt.Println()
					password = strings.TrimSpace(string(inputPassword))
					if password == "" {
						return fmt.Errorf("password cannot be empty")
					}
				}
				chiffre, err = crypt.EncryptWithPassword(plaintext, password)
			}

			if err != nil {
				return err
			}
			fmt.Println(chiffre)

			return nil
		},
	}
)

func init() {
	encryptCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "")
	encryptCmd.PersistentFlags().StringVarP(&recipientFile, "recipient-file", "r", "", "")
	encryptCmd.PersistentFlags().StringArrayVarP(&recipients, "recipients", "R", []string{}, "")
}
