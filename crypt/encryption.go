package crypt

// Encrypt encrypts a string with the given recipients (as strings) and returns the encrypted string.
func Encrypt(plaintext string, recipients []string) (string, error) {
	recipientsAge, err := parseRecipients(recipients)
	if err != nil {
		return "", err
	}

	encrypted, err := encryptString(plaintext, recipientsAge)
	if err != nil {
		return "", err
	}

	return capsule(encrypted)
}

// EncryptWithPassword encrypts a string with a password and returns the encrypted string.
func EncryptWithPassword(plaintext, password string) (string, error) {
	encrypted, err := encryptStringWithPassword(plaintext, password)
	if err != nil {
		return "", err
	}

	return capsule(encrypted)
}

// EncryptStringWithPassword encrypts a string with recipients defined in recipientFile and returns the encrypted string.
func EncryptStringWithRecipientFile(plaintext, recipientFile string) (string, error) {
	recipients, err := parseLinesFromJSONFile(recipientFile)
	if err != nil {
		return "", err
	}

	return Encrypt(plaintext, recipients)
}

// EncryptMultipleCommon encrypts multiple strings with a random secret and returns the multiple encrypted strings. The random secret gets encrypted for the recipients defined in the recipientFile (= public secret). Also returns the public secret.
func EncryptMultipleCommon(plaintexts []string, recipientFile string) ([]string, string, error) {
	privateSecret := NewPlaintextPrivateSecret()
	publicSecret, err := EncryptStringWithRecipientFile(privateSecret, recipientFile)
	if err != nil {
		return nil, "", err
	}

	results := make([]string, len(plaintexts))
	for i, plaintext := range plaintexts {
		dst, err := EncryptWithPassword(plaintext, publicSecret)
		if err != nil {
			return nil, "", err
		}
		results[i] = dst
	}

	return results, publicSecret, nil
}

func NewPlaintextPrivateSecret() string {
	return randSeq(512)
}
