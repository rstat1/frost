package crypto

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	"go.alargerobot.dev/frost/common"
)

const (
	tokenRenewTime       = 1440
	tokenRenewTimeStr    = "60s"
	vaultKVPath          = "/secret/platform/"
	vaultAppRolePrefix   = "/auth/approle/role/"
	vaultDecryptEndpoint = "/transit/decrypt/"
	vaultEncryptEndpoint = "/transit/encrypt/"
	vaultDataKeyEndpoint = "/transit/datakey/plaintext/"
)

//VaultClient ...
type VaultClient struct {
	stopRefresh         bool
	dev                 bool
	haveToken           bool
	tokenRenewerStarted bool
	client              *vault.Client
	renewer             *vault.Renewer
	currentAuthToken    string
	TokenSet            bool
	TokenSetWatch       chan bool
}

//Context extra info describing usage of a particular key. Useful in generating derived keys.
type Context map[string]string

//NewVaultClient ...
func NewVaultClient(dev bool) *VaultClient {
	c, err := vault.NewClient(&vault.Config{Address: common.CurrentConfig.VaultAddr})
	if err != nil {
		common.LogError("", err)
	}
	kms := VaultClient{client: c, dev: dev, TokenSetWatch: make(chan bool, 1), TokenSet: false}
	kms.haveToken = kms.SetAccessToken()
	return &kms
}

//GenerateKey Generates a data key. Returns plaintext and encrypted (or "sealed") forms of the key. Or an error.
func (kms *VaultClient) GenerateKey(keyID string, ctx Context) (key [32]byte, sealed []byte, e error) {
	if bytes, err := json.Marshal(&ctx); err == nil {
		payload := map[string]interface{}{
			"context": base64.StdEncoding.EncodeToString(bytes),
		}
		if newKey, err := kms.client.Logical().Write(vaultDataKeyEndpoint+keyID, payload); err == nil {
			sealedKey := newKey.Data["ciphertext"].(string)
			if notsealed, err := base64.StdEncoding.DecodeString(newKey.Data["plaintext"].(string)); err != nil {
				return key, sealed, common.LogError("generateKey", err)
			} else {
				copy(key[:], []byte(notsealed))
				return key, []byte(sealedKey), nil
			}
		} else {
			return key, sealed, common.LogError("", err)
		}
	} else {
		return key, sealed, common.LogError("", err)
	}
}

//Encrypt Encrypts a plainBlob using the key named in the global config. Returns a cipherBlob. Or an error
func (kms *VaultClient) Encrypt(plainBlob string) (string, error) {
	payload := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte(plainBlob)),
	}
	if ciphertext, err := kms.client.Logical().Write(vaultEncryptEndpoint+common.CurrentConfig.VaultKeyID, payload); err == nil {
		return ciphertext.Data["ciphertext"].(string), nil
	} else {
		return "", err
	}
}

//Decrypt Decrypts a given cipherBlob using the key named in the global config. Returns a plainBlob. Or an error
func (kms *VaultClient) Decrypt(cipherBlob string) ([]byte, error) {
	payload := map[string]interface{}{
		"ciphertext": cipherBlob,
	}
	if plaintext, err := kms.client.Logical().Write(vaultDecryptEndpoint+common.CurrentConfig.VaultKeyID, payload); err == nil {
		return base64.StdEncoding.DecodeString(plaintext.Data["plaintext"].(string))
	} else {
		return nil, err
	}
}

//WriteKeyToKVStorage Writes decryption key to the specified route in Vault
func (kms *VaultClient) WriteKeyToKVStorage(key []byte, path string) error {
	payload := map[string]interface{}{}
	payload["key"] = key
	var keyName string
	if common.DevMode {
		keyName = "-dev"
	} else {
		keyName = "-prod"
	}
	if _, e := kms.client.Logical().Write(vaultKVPath+path+keyName, payload); e != nil {
		return e
	}
	return nil
}

//ReadKeyFromKV Read a decryption key stored at the specified path from Vault
func (kms *VaultClient) ReadKeyFromKV(path string) ([]byte, error) {
	var keyName string
	if common.DevMode {
		keyName = "-dev"
	} else {
		keyName = "-prod"
	}

	if value, e := kms.client.Logical().Read(vaultKVPath + path + keyName); e != nil {
		return []byte{}, e
	} else {
		if value != nil {
			return []byte(value.Data["key"].(string)), nil
		} else {
			return []byte{}, common.LogError("", errors.New("not found"))
		}
	}
}

//UnsealKey Unseals the provided key using the provided master key. Returns plaintext key. Or an error.
func (kms *VaultClient) UnsealKey(keyID string, sealedKey []byte, ctx Context) (key [32]byte, e error) {
	if bytes, err := json.Marshal(&ctx); err == nil {
		payload := map[string]interface{}{
			"ciphertext": string(sealedKey),
			"context":    base64.StdEncoding.EncodeToString(bytes),
		}
		if unsealed, err := kms.client.Logical().Write(vaultDecryptEndpoint+keyID, payload); err == nil {
			base64Key := unsealed.Data["plaintext"].(string)
			if plainKey, err := base64.StdEncoding.DecodeString(base64Key); err == nil {
				copy(key[:], []byte(plainKey))
				return key, nil
			} else {
				return key, common.LogError("", err)
			}
		} else {
			return key, common.LogError("", err)
		}
	} else {
		return key, common.LogError("", err)
	}
}

//RenewToken Renews a token
func (kms *VaultClient) RenewToken() *vault.Secret {
	if s, e := kms.client.Auth().Token().RenewTokenAsSelf(kms.client.Token(), 28800); e == nil {
		kms.client.SetToken(s.Auth.ClientToken)
		kms.currentAuthToken = s.Auth.ClientToken
		return s
	}
	return nil
}

//SetAccessToken ...
func (kms *VaultClient) SetAccessToken() bool {
	if token, exists := os.LookupEnv("INITTOKEN"); exists {
		kms.client.SetToken(token)
		kms.TokenSet = true
		kms.TokenSetWatch <- true
	} else if common.CurrentConfig.VaultToken != "" {
		common.LogDebug("", "", "token set...")
		kms.client.SetToken(common.CurrentConfig.VaultToken)
		kms.TokenSet = true
		kms.TokenSetWatch <- true
	} else {
		common.LogWarn("", "", "no token set")
		return false
	}
	if !kms.tokenRenewerStarted {
		kms.tokenRenewalTimer()
	}
	return true
}

//GetRoleID Gets the id of the AppRole named by 'serviceName'
func (kms *VaultClient) GetRoleID(serviceName string) (string, error) {
	if kms.client.Token() == "" {
		return "", errors.New("set token first")
	}

	if value, err := kms.client.Logical().Read(vaultAppRolePrefix + "/" + serviceName + "/role-id"); err == nil {
		return value.Data["role_id"].(string), nil
	} else {
		return "", err
	}
}

//GetSecretIDAccessor Generates a short-lived single use token that can only be used to get a service's secret ID
func (kms *VaultClient) GetSecretIDAccessor() (string, error) {
	if kms.client.Token() == "" {
		return "", errors.New("set token first")
	}

	createOpts := &vault.TokenCreateRequest{
		TTL:       "5m",
		NumUses:   1,
		Renewable: common.NewFalse(),
		Policies:  []string{common.CurrentConfig.VaultFrostAppRoleName},
	}

	token, err := kms.client.Auth().Token().CreateWithRole(createOpts, common.CurrentConfig.VaultServiceAppRole)

	if err == nil {
		return token.Auth.ClientToken, nil
	} else {
		return "", err
	}
}

func (kms *VaultClient) tokenRenewalTimer() {
	go func() {
		kms.tokenRenewerStarted = true
		for {
			if !kms.stopRefresh {
				select {
				case <-time.After(kms.getTokenLeaseTime()):
					kms.RenewToken()
				}
			} else {
				break
			}
		}
	}()
}

func (kms *VaultClient) getTokenLeaseTime() time.Duration {
	var tokenDuration time.Duration
	if s, err := kms.client.Auth().Token().LookupSelf(); err == nil {
		td, _ := s.TokenTTL()
		if td > 0 {
			tokenDuration = time.Duration(td.Seconds()-20) * time.Second
		} else {
			common.LogWarn("", "", "this token has no lease duration for some reason. Defaulting to 8 hours")
			tokenDuration = time.Duration(8) * time.Hour
		}
	} else {
		common.LogError("", err)
		if !kms.dev {
			panic(err)
		} else {
			kms.stopRefresh = true
			return time.Duration(0 * time.Second)
		}
	}
	return tokenDuration
}
