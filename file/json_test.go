package file

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/mheers/inline-age/helpers"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestInitNewJSONFile(t *testing.T) {
	testFile, recipientsFile := initNewJSONFile(t)

	fmt.Println(testFile, recipientsFile)

	result, err := os.ReadFile(testFile)
	require.NoError(t, err)
	require.Contains(t, string(result), "PublicSecret")
}

func TestInitExistingJSONFile(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := path.Join(tmpDir, "test.json")
	err := os.WriteFile(testFile, []byte("{\"SomeName\": \"already given\"}"), 0600)
	require.NoError(t, err)

	recipientsFile := path.Join(tmpDir, "recipients.json")
	recipients := `{
		"__ia_config__": {
			"Recipients": {
				"marcel": "github:mheers",
				"demo": "` + demoSSHPublicKey() + `"
			}
		}
	}`

	err = os.WriteFile(recipientsFile, []byte(recipients), 0600)
	require.NoError(t, err)

	err = InitJSONFile(testFile, recipientsFile)
	require.NoError(t, err)

	fmt.Println(testFile, recipientsFile)

	result, err := os.ReadFile(testFile)
	require.NoError(t, err)
	require.Contains(t, string(result), "__ia_config__")
	require.Contains(t, string(result), "PublicSecret")
	require.Contains(t, string(result), "SomeName")
}

func TestReEncryptJSONPaths(t *testing.T) {
	testFile, recipientsFile := initNewJSONFile(t)

	err := EncryptJSONFilePath(testFile, "user.password", "juhu", demoSSHPrivateKeyFile())
	require.NoError(t, err)

	json, err := getJSONString(testFile)
	require.NoError(t, err)
	key1 := getKey(json)
	require.NotEmpty(t, key1)
	plaintextKey1, err := getPlaintextKey(json, demoSSHPrivateKeyFile())
	require.NoError(t, err)
	require.NotEmpty(t, plaintextKey1)
	encrypted1 := gjson.Get(json, "user.password").String()
	require.NotEmpty(t, encrypted1)

	plain, err := DecryptJSONFilePath(testFile, "user.password", demoSSHPrivateKeyFile())
	require.NoError(t, err)

	require.Equal(t, "juhu", plain)

	err = ReEncryptJSONFilePaths(testFile, []string{"user.password"}, demoSSHPrivateKeyFile(), recipientsFile)
	require.NoError(t, err)

	json2, err := getJSONString(testFile)
	require.NoError(t, err)
	key2 := getKey(json2)
	require.NotEmpty(t, key2)
	plaintextKey2, err := getPlaintextKey(json2, demoSSHPrivateKeyFile())
	require.NoError(t, err)
	require.NotEmpty(t, plaintextKey2)
	encrypted2 := gjson.Get(json2, "user.password").String()
	require.NotEmpty(t, encrypted2)

	plain2, err := DecryptJSONFilePath(testFile, "user.password", demoSSHPrivateKeyFile())
	require.NoError(t, err)
	require.Equal(t, "juhu", plain2)

	require.NotEqual(t, key1, key2)
	require.NotEqual(t, plaintextKey1, plaintextKey2)
	require.NotEqual(t, encrypted1, encrypted2)
}

func demoSSHPrivateKeyFile() string {
	return path.Join(getSSHTestDir(), "ssh", "id_rsa")
}

func demoSSHPubilcKeyFile() string {
	return path.Join(getSSHTestDir(), "ssh", "id_rsa.pub")
}

func demoSSHPublicKey() string {
	return string(helpers.MustReadFile(demoSSHPubilcKeyFile()))
}

func getSSHTestDir() string {
	root, err := helpers.Root()
	if err != nil {
		panic(err)
	}
	testDir := path.Join(root, "reference", "git", "test")
	return testDir
}

func initNewJSONFile(t *testing.T) (string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	testFile := path.Join(tmpDir, "test.json")
	recipientsFile := path.Join(tmpDir, "recipients.json")

	recipients := `{
		"__ia_config__": {
			"Recipients": {
				"marcel": "github:mheers",
				"demo": "` + demoSSHPublicKey() + `"
			}
		}
	}`

	err := os.WriteFile(recipientsFile, []byte(recipients), 0600)
	require.NoError(t, err)

	err = InitJSONFile(testFile, recipientsFile)
	require.NoError(t, err)

	return testFile, recipientsFile
}
