package git

import (
	"fmt"
	"path"
	"testing"

	memfs "github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCloneSSH(t *testing.T) {
	gc, err := getDemoGitSSHConnection(t)
	require.Nil(t, err)
	require.NotNil(t, gc)

	err = gc.clone()
	require.Nil(t, err)
	require.NotNil(t, gc.repo)
}

func TestCloneHTTP(t *testing.T) {
	gc, err := getDemoGitHTTPConnection(t)
	require.Nil(t, err)
	require.NotNil(t, gc)

	err = gc.clone()
	require.Nil(t, err)
	require.NotNil(t, gc.repo)
}

func TestPull(t *testing.T) {
	gc, err := getDemoGitSSHConnection(t)
	require.Nil(t, err)
	require.NotNil(t, gc)

	err = gc.clone()
	require.Nil(t, err)
	require.NotNil(t, gc.repo)

	err = gc.pull("main")
	require.NotNil(t, err)
	// require.Contains(t, err.Error(), "already up-to-date")
	require.NotNil(t, gc.repo)
}

func TestGetSSHAUth(t *testing.T) {
	auth, err := getSSHAuth(t)
	require.Nil(t, err)
	require.NotNil(t, auth)
}

func TestGetWorktree(t *testing.T) {
	gc, err := getDemoGitSSHConnection(t)
	require.Nil(t, err)
	require.NotNil(t, gc)

	wt, err := gc.getWorktree("demo")
	require.Nil(t, err)
	require.NotNil(t, wt)
}

func getSSHAuth(t *testing.T) (transport.AuthMethod, error) {
	t.Helper()
	testDir := getTestDir(t)
	return GetSSHKeyFileAuth(path.Join(testDir, "data", "id_rsa"), true)
}

func getDemoGitSSHConnection(t *testing.T) (*GitConnection, error) {
	hg, config := GetFakeGitSSHGateway(t)
	auth, err := hg.getAuth()
	if err != nil {
		return nil, err
	}

	storer := memory.NewStorage()
	fs := memfs.New()

	return NewGitConnection(auth, config.RepoURL, true, true, storer, fs, *logrus.NewEntry(logrus.StandardLogger())), nil
}

func getDemoGitHTTPConnection(t *testing.T) (*GitConnection, error) {
	hg := GetFakeGitHTTPGateway(t)
	auth, err := hg.getAuth()
	if err != nil {
		return nil, err
	}

	storer := memory.NewStorage()
	fs := memfs.New()

	return NewGitConnection(auth, hg.config.RepoURL, false, false, storer, fs, *logrus.NewEntry(logrus.StandardLogger())), nil
}

func GetFakeGitHTTPGateway(t *testing.T) *GitGateway {
	t.Helper()
	ip := getFakeGitIP(t)

	config := &Config{
		RepoURL:  fmt.Sprintf("http://%s:3000/gitea_admin/demo.git", ip),
		Username: "gitea_admin",
		Password: "admin",
	}
	gw, err := NewGitGateway(config)
	if err != nil {
		panic(err)
	}

	return gw
}
