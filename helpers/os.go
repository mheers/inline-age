package helpers

import (
	"os"
	"path"

	"github.com/adrg/xdg"
)

func PrivateKeyPath() string {
	return path.Join(xdg.Home, ".ssh", "id_rsa")
}

func MustReadFile(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return data
}
