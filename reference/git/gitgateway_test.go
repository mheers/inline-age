package git

import (
	"context"
	"fmt"
	"io"
	"path"
	"testing"
	"time"

	"github.com/mheers/inline-age/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGetSecret(t *testing.T) {
	gw, _ := GetFakeGitSSHGateway(t)

	key, value, err := gw.GetSecret("demo.json", "password")
	require.NoError(t, err)

	require.Equal(t, "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", key)
	require.Equal(t, "i2mAYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IHNjcnlwdCAweitGSEJZZWpBZEZzSkNUSk93aUFBIDE4Cnc3RzVmSzhrdys2Z05tNHI1SkxlVHVIL25vMFJkM0wya1I5eGZ4SmtLbHMKLS0tIDVmRE5zYVlmeTM4dlFHamY3STdKdzlSNXBvcTlDeEphZmlTVW43dTVIdVUKYyWc/787+BUhX0hQxr/SLBvzaugvC93T1rSPapruWszf0pYG6MKyuwjOfE6tXcWLvO3ozK5uFHrfirqaEXUD", value)
}

func GetFakeGitSSHGateway(t *testing.T) (*GitGateway, *Config) {
	t.Helper()
	ip := getFakeGitIP(t)
	testDir := getTestDir(t)

	config := &Config{
		RepoURL:    fmt.Sprintf("git@%s:gitea_admin/demo.git", ip),
		SSHKeyFile: path.Join(testDir, "ssh", "id_rsa"),
		Insecure:   true,
		Branch:     "main",
	}
	gw, err := NewGitGateway(config)
	if err != nil {
		panic(err)
	}

	return gw, config
}

func getTestDir(t *testing.T) string {
	t.Helper()
	root, err := helpers.Root()
	if err != nil {
		panic(err)
	}
	testDir := path.Join(root, "reference", "git", "test")
	return testDir
}

// getFakeGitIP returns the IP of a new (temp) fake git server
func getFakeGitIP(t *testing.T) string {
	testDir := getTestDir(t)

	compose, err := tc.NewDockerCompose(fmt.Sprintf("%s/docker-compose.yml", testDir))
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	compose.WaitForService("git", &wait.LogStrategy{
		Log:        "Starting new Web server:",
		Occurrence: 1,
	})

	container, err := compose.ServiceContainer(context.Background(), "git")
	assert.NoError(t, err, "ServiceContainer()")

	time.Sleep(6 * time.Second) // only needed to fix "database is locked (code 5 sqlite_busy)" error // TODO: find a better solution

	exitCode, r, err := container.Exec(context.Background(), []string{"sh", "/bin/init.sh"})
	require.NoError(t, err, "Exec()")
	result, err := io.ReadAll(r)
	require.NoError(t, err, "ReadAll()")
	// assert.Empty(t, string(result), "ReadAll()")
	t.Log(string(result))
	fmt.Println(string(result))
	require.Equal(t, 0, exitCode, "init could not be executed")

	ip, err := container.ContainerIP(ctx)
	assert.NoError(t, err, "ContainerIP()")

	return ip
}
