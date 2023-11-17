package reference

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mheers/inline-age/consts"
	"github.com/mheers/inline-age/crypt"
	"github.com/mheers/inline-age/file"
	"github.com/mheers/inline-age/reference/git"
	"github.com/mheers/inline-age/reference/vault"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type SecretReference struct {
	SecretStore     *SecretStore
	SecretStoreName string `json:"SecretStore"`
	GitReference    *git.Reference
	VaultReference  *vault.Reference
}

type SecretStore struct {
	Type        string `json:"Type"` // git-ssh, vault
	GitConfig   git.Config
	VaultConfig vault.Config
}

func ReadSecretStores(path string) (map[string]*SecretStore, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	jsonS := string(f)
	ssMap := gjson.Get(jsonS, consts.SecretStoresKey()).Map()

	if len(ssMap) == 0 {
		return nil, fmt.Errorf("no secret stores found")
	}

	secretStores := map[string]*SecretStore{}
	for k, v := range ssMap {
		var ss SecretStore
		err = json.Unmarshal([]byte(v.Raw), &ss)
		if err != nil {
			return nil, err
		}
		switch ss.Type {
		case "git-ssh":
			var c git.Config
			err = json.Unmarshal([]byte(v.Raw), &c)
			if err != nil {
				return nil, err
			}
			ss.GitConfig = c
		case "vault":
			var c vault.Config
			err = json.Unmarshal([]byte(v.Raw), &c)
			if err != nil {
				return nil, err
			}
			ss.VaultConfig = c
		default:
			return nil, fmt.Errorf("unknown secret store type %s", ss.Type)
		}

		secretStores[k] = &ss
	}

	return secretStores, nil
}

func ReadSecretReferences(path string, secretStores map[string]*SecretStore) (map[string]*SecretReference, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	srMap := gjson.Get(string(f), consts.SecretReferencesKey()).Map()

	if len(srMap) == 0 {
		return nil, fmt.Errorf("no secret references found")
	}

	secretReferences := map[string]*SecretReference{}
	for k, v := range srMap {
		var sr SecretReference
		err = json.Unmarshal([]byte(v.Raw), &sr)
		if err != nil {
			return nil, err
		}

		ss, ok := secretStores[sr.SecretStoreName]
		if !ok {
			return nil, fmt.Errorf("secret store %s not found", sr.SecretStoreName)
		}

		sr.SecretStore = ss
		switch sr.SecretStore.Type {
		case "git-ssh":
			var r git.Reference
			err = json.Unmarshal([]byte(v.Raw), &r)
			if err != nil {
				return nil, err
			}
			r.Config = &sr.SecretStore.GitConfig
			sr.GitReference = &r
		case "vault":
			var r vault.Reference
			err = json.Unmarshal([]byte(v.Raw), &r)
			if err != nil {
				return nil, err
			}
			r.Config = &sr.SecretStore.VaultConfig
			sr.VaultReference = &r
		default:
			return nil, fmt.Errorf("unknown secret store type %s", sr.SecretStore.Type)
		}

		secretReferences[k] = &sr
	}

	return secretReferences, nil
}

func (r *SecretReference) getEncryptedValue() (string, string, error) {
	switch r.SecretStore.Type {
	case "git-ssh":
		key, secret, err := r.GitReference.GetSecret()
		if err != nil {
			return "", "", err
		}
		return key, secret, nil
	case "vault":
		secret, err := r.VaultReference.GetSecret()
		if err != nil {
			return "", "", err
		}
		return "", secret, nil

	default:
		return "", "", nil
	}
}

func (r *SecretReference) PlaintTextValue(identityFile string) (string, error) {
	key, data, err := r.getEncryptedValue()
	if err != nil {
		return "", err
	}

	if key == "" { // no key, so it's not encrypted
		return data, nil
	}

	keyDec, err := crypt.DecryptFromIdentityFile(key, identityFile)
	if err != nil {
		return "", err
	}
	return crypt.DecryptFromPassword(data, keyDec)
}

func ResolveReferences(jsonFile, refsFile, identityFile string) error {
	secretStores, err := ReadSecretStores(refsFile)
	if err != nil {
		return err
	}

	refs, err := ReadSecretReferences(refsFile, secretStores)
	if err != nil {
		return err
	}

	j, err := os.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	json := string(j)

	key := gjson.Get(json, consts.PublicSecretKey()).String()
	keyDec, err := crypt.DecryptFromIdentityFile(key, identityFile)
	if err != nil {
		return err
	}

	referenceMapping, err := getReferenceMapping(json)
	if err != nil {
		return err
	}

	for m, p := range referenceMapping {
		r, err := refs[p].PlaintTextValue(identityFile)
		if err != nil {
			return err
		}

		reEncrypted, err := crypt.EncryptWithPassword(r, keyDec)
		if err != nil {
			return err
		}

		json, err = sjson.Set(json, m, reEncrypted)
		if err != nil {
			return err
		}

		json, err = file.EnsurePathInSecretPaths(json, m)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(jsonFile, []byte(json), 0600)
	if err != nil {
		return err
	}

	return nil
}

func getReferenceMapping(json string) (map[string]string, error) {
	referenceMapping := gjson.Get(json, consts.PathReferenceMappingKey()).Map()
	if len(referenceMapping) == 0 {
		return nil, fmt.Errorf("no referenceMapping found")
	}

	result := make(map[string]string)
	for m, p := range referenceMapping {
		result[m] = p.Str
	}

	return result, nil
}
