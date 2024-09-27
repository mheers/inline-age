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
	decryptPCmd = &cobra.Command{
		Use:     "decryptp [chiffre]",
		Short:   "decrypts a string with a password",
		Long:    ``,
		Aliases: []string{"dp"},
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
			if password == "" {
				fmt.Fprint(os.Stderr, "Enter password: ")
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
				fmt.Fprint(os.Stderr)
				password = strings.TrimSpace(string(inputPassword))
				if password == "" {
					return fmt.Errorf("password cannot be empty")
				}
			}

			plaintext, err := crypt.DecryptFromPassword(chiffre, password)
			if err != nil {
				return err
			}
			fmt.Println(plaintext)

			return nil
		},
	}
)

func init() {
	decryptPCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "")
}
