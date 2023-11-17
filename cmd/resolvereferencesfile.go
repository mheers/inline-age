package cmd

import (
	"fmt"

	"github.com/mheers/inline-age/helpers"
	"github.com/mheers/inline-age/reference"
	"github.com/spf13/cobra"
)

var (
	referenceFile        string
	referenceFileDefault string = ""

	resolveReferencesFileCmd = &cobra.Command{
		Use:     "resolve-reference-file [file]",
		Short:   "resolve-references a a file",
		Long:    ``,
		Aliases: []string{"rrf"},
		RunE: func(cmd *cobra.Command, args []string) error {
			helpers.SetLogLevel(logLevelFlag)

			if len(args) == 0 {
				return fmt.Errorf("not enough arguments")
			}

			fileName := args[0]

			if referenceFile == "" {
				cmd.Println("no reference file set, using references from secret file")
				referenceFile = fileName
			}

			err := reference.ResolveReferences(fileName, referenceFile, identityFile)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	resolveReferencesFileCmd.PersistentFlags().StringVarP(&referenceFile, "reference-file", "r", referenceFileDefault, "if not set the references defined in the secret file will be used")
	resolveReferencesFileCmd.PersistentFlags().StringVarP(&identityFile, "identity-file", "i", helpers.PrivateKeyPath(), "")
}
