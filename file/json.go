package file

import (
	"errors"
	"fmt"
	"os"

	"github.com/mheers/inline-age/consts"
	"github.com/mheers/inline-age/crypt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
)

func EncryptJSONFilePath(file, path, plaintext, identityFile string) error {
	json, err := getJSONString(file)
	if err != nil {
		return err
	}

	keyDec, err := getPlaintextKey(json, identityFile)
	if err != nil {
		return err
	}

	json, err = setPlaintextWithKey(json, path, plaintext, keyDec)
	if err != nil {
		return err
	}

	json, err = EnsurePathInSecretPaths(json, path)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, []byte(json), 0600)
	if err != nil {
		return err
	}

	return nil
}

func EnsurePathInSecretPaths(json, path string) (string, error) {
	paths := []string{}
	pathsG := gjson.Get(json, consts.SecretPathsKey()).Array()

	for _, path := range pathsG {
		paths = append(paths, path.String())
	}

	var err error
	if !slices.Contains(paths, path) {
		paths = append(paths, path)
		json, err = sjson.Set(json, consts.SecretPathsKey(), paths)
	}

	return json, err
}

func setPlaintextWithKey(json, path, plaintext, keyDec string) (string, error) {
	encrypted, err := crypt.EncryptWithPassword(plaintext, keyDec)
	if err != nil {
		return "", err
	}

	json, err = sjson.Set(json, path, encrypted)
	if err != nil {
		return "", err
	}

	return json, nil
}

func DecryptJSONFilePath(file, path, identityFile string) (string, error) {
	json, err := getJSONString(file)
	if err != nil {
		return "", err
	}

	return getPlaintext(json, path, identityFile)
}

func getJSONString(file string) (string, error) {
	j, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	json := string(j)

	return json, nil
}

func getKey(json string) string {
	return gjson.Get(json, consts.PublicSecretKey()).String()
}

func getPlaintextKey(json, identityFile string) (string, error) {
	key := getKey(json)
	keyDec, err := crypt.DecryptFromIdentityFile(key, identityFile)
	if err != nil {
		return "", err
	}
	return keyDec, nil
}

func getPlaintext(json, path, identityFile string) (string, error) {
	keyDec, err := getPlaintextKey(json, identityFile)
	if err != nil {
		return "", err
	}

	return getPlaintextFromKey(json, path, keyDec)
}

func getPlaintextFromKey(json, path, keyDec string) (string, error) {
	chiffre := gjson.Get(json, path).String()

	plaintext, err := crypt.DecryptFromPassword(chiffre, keyDec)
	if err != nil {
		return "", err
	}
	return plaintext, nil
}

func ReencryptJSONFile(file string, identityFile, recipientFile string) error {
	json, err := getJSONString(file)
	if err != nil {
		return err
	}

	paths := []string{}
	pathsG := gjson.Get(json, consts.SecretPathsKey()).Array()

	for _, path := range pathsG {
		paths = append(paths, path.String())
	}

	return ReEncryptJSONFilePaths(file, paths, identityFile, recipientFile)
}

func ReEncryptJSONFilePaths(file string, paths []string, identityFile, recipientFile string) error {
	json, err := getJSONString(file)
	if err != nil {
		return err
	}

	keyDec, err := getPlaintextKey(json, identityFile)
	if err != nil {
		return err
	}

	plaintexts := make(map[string]string)

	// get all old plaintexts
	for _, path := range paths {
		plaintext, err := getPlaintextFromKey(json, path, keyDec)
		if err != nil {
			return err
		}
		plaintexts[path] = plaintext
	}

	// create new public secret
	newPlaintextPublicSecret := crypt.NewPlaintextPrivateSecret()
	publicSecret, err := crypt.EncryptStringWithRecipientFile(newPlaintextPublicSecret, recipientFile)
	if err != nil {
		return err
	}

	// reencrypt all plaintexts with new public secret
	for path, plaintext := range plaintexts {
		json, err = setPlaintextWithKey(json, path, plaintext, newPlaintextPublicSecret)
		if err != nil {
			return err
		}
	}

	// set new public secret
	json, err = sjson.Set(json, consts.PublicSecretKey(), publicSecret)
	if err != nil {
		return err
	}

	// write back
	err = os.WriteFile(file, []byte(json), 0600)
	if err != nil {
		return err
	}

	return nil
}

func InitJSONFile(file, recipientFile string) error {
	_, publicSecret, err := crypt.EncryptMultipleCommon([]string{}, recipientFile)
	if err != nil {
		return err
	}

	var json string

	finfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist, will be created")
		} else {
			return err
		}
	} else if finfo.IsDir() {
		return errors.New("file is dir")
	} else {
		j, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		json = string(j)
	}

	json, err = sjson.Set(json, consts.PublicSecretKey(), publicSecret)
	if err != nil {
		return err
	}
	json, err = sjson.Set(json, consts.SecretPathsKey(), []string{})
	if err != nil {
		return err
	}
	json, err = sjson.Set(json, consts.SecretStoresKey(), map[string]string{})
	if err != nil {
		return err
	}
	json, err = sjson.Set(json, consts.SecretReferencesKey(), map[string]string{})
	if err != nil {
		return err
	}
	json, err = sjson.Set(json, consts.PathReferenceMappingKey(), map[string]string{})
	if err != nil {
		return err
	}

	recipientsMap, err := crypt.GetRecipientsFromJSONFile(recipientFile)
	if err != nil {
		return err
	}
	json, err = sjson.Set(json, consts.RecipientsKey(), recipientsMap)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, []byte(json), 0600)
	if err != nil {
		return err
	}

	return nil
}
