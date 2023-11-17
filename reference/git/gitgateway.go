package git

import (
	"fmt"
	"os"
	"path"

	osfs "github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mheers/inline-age/consts"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const defaultBranch = "main"

type Config struct {
	RepoURL    string `json:"RepoURL" yaml:"RepoURL"`
	Branch     string `json:"Branch" yaml:"Branch"`
	Username   string `json:"Username" yaml:"Username"`
	Password   string `json:"Password" yaml:"Password"`
	Insecure   bool   `json:"Insecure" yaml:"Insecure"`
	SSHKey     string `json:"SshKey" yaml:"SshKey"`
	SSHKeyFile string `json:"SshKeyFile" yaml:"SshKeyFile"`
}

type Reference struct {
	Config   *Config
	JSONFile string `json:"JsonFile" yaml:"JsonFile"`
	Path     string `json:"Path" yaml:"Path"`
	// YAMLFile   string `json:"yamlFile" yaml:"yamlFile"` // TODO: implement
	// CUEFile   string `json:"cueFile" yaml:"cueFile"` // TODO: implement
}

type GitGateway struct {
	config *Config
}

func NewGitGateway(cfg *Config) (*GitGateway, error) {
	if cfg.SSHKeyFile != "" && cfg.SSHKey != "" {
		return nil, fmt.Errorf("cannot specify both sshKey and sshKeyFile")
	}

	if (cfg.SSHKeyFile != "" || cfg.SSHKey != "") && cfg.Password != "" {
		return nil, fmt.Errorf("cannot specify both sshKey and username")
	}

	if cfg.Branch == "" {
		fmt.Printf("no branch specified, using default branch: %s\n", defaultBranch)
		cfg.Branch = defaultBranch
	}

	gw := &GitGateway{
		config: cfg,
	}

	return gw, nil
}

func (gw *GitGateway) GetSecret(filePath, secretPath string) (string, string, error) {
	// create a temp dir where we clone the repo
	tmpDir, err := os.MkdirTemp("", "git-src")
	if err != nil {
		return "", "", err
	}

	storer := memory.NewStorage()

	// create a filesystem backed by the temp dir
	fs := osfs.New(tmpDir)

	// authenticate
	auth, err := gw.getAuth()
	if err != nil {
		return "", "", err
	}

	// make a connection to the git repo
	conn := NewGitConnection(auth, gw.config.RepoURL, gw.config.Insecure, gw.config.Insecure, storer, fs, *logrus.NewEntry(logrus.StandardLogger()))

	// clone the repo
	if err = conn.clone(); err != nil {
		return "", "", err
	}

	// checkout the branch
	_, err = conn.GetWorktree(gw.config.Branch)
	if err != nil {
		return "", "", err
	}

	data, err := os.ReadFile(path.Join(tmpDir, filePath))
	if err != nil {
		return "", "", err
	}

	jsonS := string(data)

	key := gjson.Get(jsonS, consts.PublicSecretKey()).String()
	secret := gjson.Get(jsonS, secretPath).String()

	return key, secret, nil
}

func (gw *GitGateway) getAuth() (transport.AuthMethod, error) {
	if gw.config.Username != "" && gw.config.Password != "" {
		return GetBasicAuth(gw.config.Username, gw.config.Password)
	}

	if gw.config.SSHKey != "" {
		return GetSSHKeyAuth([]byte(gw.config.SSHKey), gw.config.Insecure)
	}

	if gw.config.SSHKeyFile != "" {
		return GetSSHKeyFileAuth(gw.config.SSHKeyFile, gw.config.Insecure)
	}

	return nil, nil
}

func (r *Reference) GetSecret() (string, string, error) {
	gitGW, err := NewGitGateway(r.Config)
	if err != nil {
		return "", "", err
	}
	key, secret, err := gitGW.GetSecret(r.JSONFile, r.Path)
	if err != nil {
		return "", "", err
	}

	return key, secret, nil
}
