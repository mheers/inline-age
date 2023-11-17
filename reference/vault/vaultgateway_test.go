package vault

import (
	"testing"

	"context"
	"fmt"
	"io"
	"path"

	"github.com/mheers/inline-age/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestAll(t *testing.T) {
	gw := GetDemoVaultGateway(t)
	require.NotNil(t, gw)

	client, err := gw.vaultClient()
	require.NoError(t, err)
	require.NotNil(t, client)

	// read a secret from vault
	secret, err := gw.GetSecret("kv/foo", "bar")

	require.NoError(t, err)
	require.NotNil(t, secret)
	require.Equal(t, "baz", secret)
}

func getTestDir() string {
	root, err := helpers.Root()
	if err != nil {
		panic(err)
	}
	testDir := path.Join(root, "reference", "git", "test")
	return testDir
}

func GetDemoVaultGateway(t *testing.T) *VaultGateway {
	testDir := getTestDir()

	compose, err := tc.NewDockerCompose(fmt.Sprintf("%s/docker-compose.yml", testDir))
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	compose.WaitForService("vault", &wait.LogStrategy{
		Log:        "Development mode should NOT be used in production installations!",
		Occurrence: 1,
	})

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	container, err := compose.ServiceContainer(context.Background(), "vault")
	assert.NoError(t, err, "ServiceContainer()")

	exitCode, r, err := container.Exec(context.Background(), []string{"sh", "/bin/init.sh"})
	require.NoError(t, err, "Exec()")
	require.Equal(t, 0, exitCode, "init could not be executed")
	result, err := io.ReadAll(r)
	require.NoError(t, err, "ReadAll()")
	fmt.Println(string(result))

	ip, err := container.ContainerIP(context.Background())
	assert.NoError(t, err, "ContainerIP()")

	cfg := &Config{
		Endpoint: "http://" + ip + ":8200",
		Token:    "root",
		Insecure: true,
	}

	return NewVaultGateway(cfg)
}
