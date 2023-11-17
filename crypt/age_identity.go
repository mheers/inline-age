package crypt

import (
	"fmt"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
)

// parseSSHIdentityFromIdentityFile parses an SSH identity file and returns a slice of age.Identity.
func parseSSHIdentityFromIdentityFile(identityFile string) ([]age.Identity, error) {
	privateKey, err := privateKey(identityFile)
	if err != nil {
		return nil, err
	}

	return parseSSHIdentity(identityFile, privateKey)
}

// parseSSHIdentity parses an SSH identity (as byte slice) and returns a slice of age.Identity.
func parseSSHIdentity(name string, pemBytes []byte) ([]age.Identity, error) {
	id, err := agessh.ParseIdentity(pemBytes)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok {
		pubKey := sshErr.PublicKey
		if pubKey == nil {
			pubKey, err = readPubFile(name)
			if err != nil {
				return nil, err
			}
		}
		passphrasePrompt := func() ([]byte, error) {
			pass, err := readSecret(fmt.Sprintf("Enter passphrase for %q:", name))
			if err != nil {
				return nil, fmt.Errorf("could not read passphrase for %q: %v", name, err)
			}
			return pass, nil
		}
		i, err := agessh.NewEncryptedSSHIdentity(pubKey, pemBytes, passphrasePrompt)
		if err != nil {
			return nil, err
		}
		return []age.Identity{i}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("malformed SSH identity in %q: %v", name, err)
	}

	return []age.Identity{id}, nil
}
