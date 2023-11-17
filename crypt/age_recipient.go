package crypt

import (
	"fmt"
	"os"
	"strings"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"github.com/mheers/inline-age/consts"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// parseLinesFromTextFile parses a text file into a slice of strings. It ignores lines that are empty or start with "#".
func parseLinesFromTextFile(filePath string) ([]string, error) {
	if filePath == "" {
		return nil, fmt.Errorf("recipient file path is empty")
	}
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	filtered := []string{}
	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}
		filtered = append(filtered, line)
	}

	return filtered, nil
}

// GetRecipientsFromJSONFile parses a text file into a slice of strings. It ignores keys that are empty or start with "#".
func GetRecipientsFromJSONFile(filePath string) (map[string]string, error) {
	if filePath == "" {
		return nil, fmt.Errorf("recipient file path is empty")
	}
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json := string(b)

	recipients := map[string]string{}
	recipientsG := gjson.Get(json, consts.RecipientsKey()).Map()

	for name, recipientG := range recipientsG {
		logrus.Debug("Found recipient: ", name)
		recipient := recipientG.String()
		if recipient == "" {
			continue
		}

		if strings.HasPrefix(recipient, "#") {
			continue
		}
		recipients[name] = recipient
	}

	return recipients, nil
}

// parseLinesFromJSONFile parses a text file into a slice of strings. It ignores keys that are empty or start with "#".
func parseLinesFromJSONFile(filePath string) ([]string, error) {
	recipientsMap, err := GetRecipientsFromJSONFile(filePath)
	if err != nil {
		return nil, err
	}

	recipients := []string{}
	for _, recipient := range recipientsMap {
		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

// parseRecipients parses a slice of strings into a slice of age.Recipient.
// It supports the following recipient types:
// - ssh keys (e.g. "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@host")
// - github usernames (e.g. "github:mheers")
// - gitlab usernames (e.g. "gitlab:mheers")
// Returns an error if any of the recipients is invalid.
func parseRecipients(recipientLines []string) ([]age.Recipient, error) {
	var recipients []age.Recipient
	for _, line := range recipientLines {
		r, err := parseRecipient(line)
		if err != nil {
			logrus.Errorf("could not parse recipient %q: %v", line, err)
			return nil, err
		}
		recipients = append(recipients, r)
	}
	return recipients, nil
}

// ParseRecipients parses a slice of strings into a slice of age.Recipient.
// It supports the following recipient types:
// - ssh keys (e.g. "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@host")
// - github usernames (e.g. "github:mheers")
// - gitlab usernames (e.g. "gitlab:mheers")
func parseRecipient(arg string) (age.Recipient, error) {
	switch {
	case strings.HasPrefix(arg, "ssh-"):
		return agessh.ParseRecipient(arg)
	case strings.HasPrefix(arg, "github:"):
		name := strings.TrimPrefix(arg, "github:")
		url := fmt.Sprintf("https://github.com/%s.keys", name)
		body, err := download(url)
		if err != nil {
			return nil, err
		}
		return agessh.ParseRecipient(body)
	case strings.HasPrefix(arg, "gitlab:"):
		name := strings.TrimPrefix(arg, "gitlab:")
		url := fmt.Sprintf("https://gitlab.com/%s.keys", name)
		body, err := download(url)
		if err != nil {
			return nil, err
		}
		return agessh.ParseRecipient(body)
	}

	return nil, fmt.Errorf("unknown recipient type: %q", arg)
}
