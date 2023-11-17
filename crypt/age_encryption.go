package crypt

import (
	"bytes"
	"io"

	"filippo.io/age"
)

// encryptString encrypts a string with the given recipients and returns the encrypted bytes.
func encryptString(in string, recipients []age.Recipient) ([]byte, error) {
	out := &bytes.Buffer{}
	w, err := age.Encrypt(out, recipients...)
	if err != nil {
		return nil, err
	}

	if _, err := io.WriteString(w, in); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// encryptStringWithPassword encrypts a string with a password and returns the encrypted string.
func encryptStringWithPassword(in, password string) ([]byte, error) {
	out := &bytes.Buffer{}
	recipient, err := age.NewScryptRecipient(password)
	if err != nil {
		return nil, err
	}
	w, err := age.Encrypt(out, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := io.WriteString(w, in); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// decryptToString decrypts a byte slice with the given age identities and returns the decrypted string.
func decryptToString(in []byte, identities []age.Identity) (string, error) {
	r, err := age.Decrypt(bytes.NewBuffer(in), identities...)
	if err != nil {
		return "", err
	}

	responseB, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(responseB), nil
}

// decryptToStringWithPassword decrypts a byte slice with a password and returns the decrypted string.
func decryptToStringWithPassword(in []byte, password string) (string, error) {
	identity, err := age.NewScryptIdentity(password)
	if err != nil {
		return "", err
	}
	r, err := age.Decrypt(bytes.NewBuffer(in), identity)
	if err != nil {
		return "", err
	}

	responseB, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(responseB), nil
}
