package reference

import (
	"os"
	"path"
	"testing"

	"github.com/mheers/inline-age/helpers"
	"github.com/stretchr/testify/require"
)

func TestReadSecretStore(t *testing.T) {
	demoDataFile, referencesFile := createDemoData(t)
	require.NotEmpty(t, demoDataFile)
	require.NotEmpty(t, referencesFile)

	r, err := ReadSecretStores(referencesFile)
	require.NoError(t, err)
	require.Equal(t, "git-ssh", r["git"].Type)
}

func TestReference(t *testing.T) {
	demoDataFile, referencesFile := createDemoData(t)
	require.NotEmpty(t, demoDataFile)
	require.NotEmpty(t, referencesFile)

	ss, err := ReadSecretStores(referencesFile)
	require.NoError(t, err)

	r, err := ReadSecretReferences(referencesFile, ss)
	require.NoError(t, err)

	require.Equal(t, "git-ssh", r["password"].SecretStore.Type)
}

func createDemoData(t *testing.T) (string, string) {
	t.Helper()

	dir := t.TempDir()

	referencePath := path.Join(dir, "references.json")
	referenceData := []byte(`{
		"__ia_config__": {
			"SecretStores": {
				"git": {
					"Type": "git-ssh",
					"GitConfig": {
						"RepoURL": "",
						"SshKeyFile": "",
						"JsonFile": "demo.json",
						"Path": "password",
						"Branch": "main"
					}
				}
			},
			"SecretReferences": {
				"password": {
					"SecretStore": "git",
					"file": "demo.json",
					"path": "password"
				}
			}
		}
	}`)
	err := os.WriteFile(referencePath, referenceData, 0644)
	require.NoError(t, err)

	demoPath := path.Join(dir, "demos.json")
	demoData := []byte(`{
		"name": "will be a secret name",
		"__ia_config__": {
			"PublicSecret": "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",
			"PathReferenceMapping": {
				"user.password": "password"
			}
		}
	  }`)
	err = os.WriteFile(demoPath, demoData, 0644)
	require.NoError(t, err)

	return demoPath, referencePath
}

// TODO: mock git or vault here
func TestResolveReferences(t *testing.T) {
	t.Skip()
	demoDataFile, referenceDataFile := createDemoData(t)

	err := ResolveReferences(demoDataFile, referenceDataFile, helpers.PrivateKeyPath())
	require.NoError(t, err)
}
