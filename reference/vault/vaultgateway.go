package vault

import (
	"errors"

	"github.com/hashicorp/vault/api"
)

type Config struct {
	Endpoint string
	Token    string
	Insecure bool
}

type Reference struct {
	Config *Config
	Name   string
	Key    string
}

func ReadConfig() (*Config, error) {
	c := &Config{}

	return c, nil
}

type VaultGateway struct {
	config *Config
}

func NewVaultGateway(config *Config) *VaultGateway {
	gw := &VaultGateway{
		config: config,
	}

	return gw
}

func (gw *VaultGateway) GetSecret(name, key string) (string, error) {
	client, err := gw.vaultClient()
	if err != nil {
		panic(err)
	}

	secretData, err := client.Logical().Read(name)
	if err != nil {
		return "", err
	}

	if secretData == nil {
		return "", errors.New("secret not found")
	}

	secret := secretData.Data

	if secret == nil {
		return "", errors.New("secret not found")
	}

	if secret[key] == nil {
		return "", errors.New("key not found")
	}

	return secret[key].(string), nil
}

func (gw *VaultGateway) vaultClient() (*api.Client, error) {
	return GetVaultClient(gw.config.Endpoint, gw.config.Token, gw.config.Insecure)
}

func GetVaultClient(endpoint, token string, insecure bool) (*api.Client, error) {
	var vclient *api.Client
	conf := api.DefaultConfig()
	var err error
	if insecure {
		conf.ConfigureTLS(&api.TLSConfig{Insecure: insecure})
	}
	vclient, err = api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	vclient.SetAddress(endpoint)

	if token == "" {
		return nil, errors.New("token for Vault required")
	}

	vclient.SetToken(token)
	return vclient, nil
}

func (r *Reference) GetSecret() (string, error) {
	return NewVaultGateway(r.Config).GetSecret(r.Name, r.Key)
}
