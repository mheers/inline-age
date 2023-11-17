package crypt

import "filippo.io/age"

// Decrypt decrypts a string with the given age identities and returns the decrypted string.
func Decrypt(ciffre string, identityAge []age.Identity) (string, error) {
	uncompressed, err := uncapsule(ciffre)
	if err != nil {
		return "", err
	}

	plaintext, err := decryptToString(uncompressed, identityAge)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}

// DecryptFromPassword decrypts a string with a password and returns the decrypted string.
func DecryptFromPassword(ciffre, password string) (string, error) {
	uncompressed, err := uncapsule(ciffre)
	if err != nil {
		return "", err
	}

	plaintext, err := decryptToStringWithPassword(uncompressed, password)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}

// DecryptFromIdentityFile decrypts a string with the given identity file and returns the decrypted string.
func DecryptFromIdentityFile(ciffre, identityFile string) (string, error) {
	identity, err := parseSSHIdentityFromIdentityFile(identityFile)
	if err != nil {
		return "", err
	}

	uncompressed, err := uncapsule(ciffre)
	if err != nil {
		return "", err
	}

	plaintext, err := decryptToString(uncompressed, identity)
	if err != nil {
		return "", err
	}

	return plaintext, nil
}

// DecryptMultipleCommon decrypts multiple strings with a public secret (chiffre) and returns the multiple decrypted strings.
func DecryptMultipleCommon(enctexts []string, chiffre, identityFile string) ([]string, error) {
	privateSecret, err := DecryptFromIdentityFile(chiffre, identityFile)
	if err != nil {
		return nil, err
	}

	results := make([]string, len(enctexts))
	for i, enctext := range enctexts {
		uncompressed, err := uncapsule(enctext)
		if err != nil {
			return nil, err
		}

		dst, err := decryptToStringWithPassword(uncompressed, privateSecret)
		if err != nil {
			return nil, err
		}

		results[i] = dst
	}

	return results, nil
}
