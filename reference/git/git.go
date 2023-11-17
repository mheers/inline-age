package git

import (
	"errors"
	"io"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	go_git_ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sirupsen/logrus"

	billy "github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"

	"golang.org/x/crypto/ssh"
)

type GitConnection struct {
	auth                     transport.AuthMethod
	repoURL                  string
	fs                       billy.Filesystem
	storer                   *memory.Storage
	progress                 sideband.Progress
	repo                     *git.Repository
	remote                   *git.Remote
	remoteName               string
	insecureSkipTLS          bool
	insecureIgnoreUnkownHost bool
	logger                   logrus.Entry
}

/*
NewGitConnection creates and returns a new git connection
*/
func NewGitConnection(auth transport.AuthMethod, repoURL string, insecureSkipTLS, insecureIgnoreUnkownHost bool, storer *memory.Storage, fs billy.Filesystem, logger logrus.Entry) *GitConnection {
	remoteName := "origin"

	if fs == nil {
		fs = memfs.New()
	}

	return &GitConnection{
		auth:                     auth,
		repoURL:                  repoURL,
		remoteName:               remoteName,
		fs:                       fs,
		storer:                   storer,
		insecureSkipTLS:          insecureSkipTLS,
		insecureIgnoreUnkownHost: insecureIgnoreUnkownHost,
		logger:                   logger,
	}
}

func (gc *GitConnection) GetFileContent(branch, path string) ([]byte, error) {
	bf, err := gc.GetFile(branch, path)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(bf)
}

func (gc *GitConnection) GetFile(branch, path string) (billy.File, error) {
	w, err := gc.GetWorktree(branch)
	if err != nil {
		return nil, err
	}

	return w.Filesystem.Open(path)
}

func (gc *GitConnection) GetBranches() ([]string, error) {
	if err := gc.initRepo(); err != nil {
		return nil, err
	}

	rem, err := gc.repo.Remote(gc.remoteName)
	if err != nil {
		return nil, err
	}
	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filters the references list and only keeps branches
	var branches []string
	for _, ref := range refs {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
	}
	return branches, nil
}

func (gc *GitConnection) clone() error {
	path := gc.fs.Root()
	repo, err := git.Clone(gc.storer, gc.fs, &git.CloneOptions{
		URL:             gc.repoURL,
		Progress:        gc.progress,
		Auth:            gc.auth,
		InsecureSkipTLS: gc.insecureSkipTLS,

		Depth: 1,
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			repo, err = git.PlainOpen(path)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	gc.repo = repo
	return nil
}

func (gc *GitConnection) init() error {
	err := gc.initRepo()
	if err != nil {
		return err
	}
	err = gc.initRemote()
	if err != nil {
		return err
	}
	return nil
}

func (gc *GitConnection) initRepo() error {
	// clone if not yet cloned
	if gc.repo == nil {
		return gc.clone()
	}
	return nil
}

func (gc *GitConnection) initRemote() error {
	if gc.remote != nil {
		return nil
	}
	remotes, err := gc.repo.Remotes()
	if err != nil {
		return err
	}
	if len(remotes) == 0 {
		return errors.New("no remotes found")
	}
	gc.logger.Debugf("git %s: initRemotes: ", remotes)
	gc.remote = remotes[0]
	return nil
}

func (gc *GitConnection) pull(referenceString string) error {

	if gc.repo == nil {
		return errors.New("pull: repo may not be nil")
	}
	w, err := gc.repo.Worktree()
	if err != nil {
		return err
	}

	reference := plumbing.NewBranchReferenceName(referenceString)

	return w.Pull(&git.PullOptions{
		RemoteName:      gc.remoteName,
		Progress:        gc.progress,
		Auth:            gc.auth,
		InsecureSkipTLS: gc.insecureSkipTLS,
		ReferenceName:   reference,
	})
}

func (gc *GitConnection) GetWorktree(branch string) (*git.Worktree, error) {
	return gc.getWorktree(branch)
}

func (gc *GitConnection) getWorktree(reference string) (*git.Worktree, error) {
	if err := gc.init(); err != nil {
		return nil, err
	}

	w, err := gc.repo.Worktree()
	if err != nil {
		return nil, err
	}

	// try to find the remote branch
	err = gc.repo.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Auth:     gc.auth,
	})
	if err != nil {
		// ignore "already up to date" errors
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			gc.logger.Println("ignoring 'already up to date'")
		} else {
			return nil, err
		}
	}

	h, _ := gc.repo.ResolveRevision(plumbing.Revision(reference))

	// checkout the wanted branch
	err = w.Checkout(&git.CheckoutOptions{
		Hash:  *h,
		Force: true,
		Keep:  true,
	})
	if err != nil {
		return nil, err
	}

	// // always pull the latest commits
	// err = gc.pull(branch)
	// if err != nil {
	// 	// ignore "already up to date" errors
	// 	if errors.Is(err, git.NoErrAlreadyUpToDate) {
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	return w, nil
}

func GetSSHKeyFileAuth(privateSshKeyFile string, insecureIgnoreHostKey bool) (transport.AuthMethod, error) {
	sshKey, err := os.ReadFile(privateSshKeyFile)
	if err != nil {
		return nil, err
	}
	return GetSSHKeyAuth(sshKey, insecureIgnoreHostKey)
}

func GetSSHKeyAuth(privateSshKey []byte, insecureIgnoreHostKey bool) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	signer, err := ssh.ParsePrivateKey(privateSshKey)
	if err != nil {
		return nil, err
	}
	auth = &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	auth.(*go_git_ssh.PublicKeys).HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return auth, nil
}

func GetBasicAuth(username, password string) (transport.AuthMethod, error) {
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}, nil
}
