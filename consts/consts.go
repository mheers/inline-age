package consts

import "fmt"

const (
	MainKey = "__ia_config__" // defines the main config entrypoint in a json secret file
)

// PublicSecretKey: config key under which the encrypted secret is stored with which the other secrets are encrypted
func PublicSecretKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "PublicSecret")
}

// SecretPathsKey: config key under which the absolute json paths to all secrets are stored - needed for reencryptoin of all secrets
func SecretPathsKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "SecretPaths")
}

// SecretStoresKey: config key under which the secret stores with their own credentials are configured
func SecretStoresKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "SecretStores")
}

// SecretReferencesKey: config key under which secrets from secret stores can be referenced
func SecretReferencesKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "SecretReferences")
}

// PathReferenceMappingKey: config key under which the mapping from absolute json paths to references is stored
func PathReferenceMappingKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "PathReferenceMapping")
}

// RecipientsKey: config key under which the recipients of this file are stored; these recipients can decrypt the PublicSecret
func RecipientsKey() string {
	return fmt.Sprintf("%s.%s", MainKey, "Recipients")
}
