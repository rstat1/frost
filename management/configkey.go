// Borrowed much of this file from the similarly named file in the Minio project, which is licensed under the
// Apache 2.0 license.

package management

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"git.m/devapp/common"
	"github.com/minio/sio"
)

//KeyInfo ...
type KeyInfo struct {
	IV        string `json:"iv"`
	SealedKey string `json:"key"`
	BlobName  string `json:"name"`
}

//ConfigEncryptionKey ...
type ConfigEncryptionKey struct {
	SealedMasterKey []byte  `json:"masterKey"`
	EntryKey        KeyInfo `json:"valueKey"`
}

//Key ...
type Key [32]byte

//GenerateKey ...
func GenerateKey(masterkey []byte, file string) Key {
	var key [32]byte
	var nonce [32]byte

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		common.LogError("", err)
		fmt.Printf("Failed to read random data: %v", err) // add error handling
	}

	sha := sha256.New()
	sha.Write(masterkey[:])
	sha.Write(nonce[:])
	sha.Sum(key[:0])

	return key
}

//Seal ...
func (k *Key) Seal(master []byte, path string) (sealedKey KeyInfo, e error) {
	var encryptedKey bytes.Buffer

	iv := k.generateIV()

	mac := hmac.New(sha256.New, master)
	mac.Write(iv[:])
	mac.Write([]byte(path))
	mac.Write([]byte("FROST-HMAC-SHA256"))

	config := sio.Config{MinVersion: sio.Version20, Key: mac.Sum(nil)}
	if num, err := sio.Encrypt(&encryptedKey, bytes.NewReader(k[:]), config); num != 64 || err != nil {
		return KeyInfo{}, common.LogError("", err)
	}
	sealedKey.IV = base64.StdEncoding.EncodeToString(iv[:])
	sealedKey.BlobName = path
	sealedKey.SealedKey = base64.StdEncoding.EncodeToString(encryptedKey.Bytes())
	return sealedKey, nil
}

//Unseal ...
func (k *Key) Unseal(master []byte, sealedKey KeyInfo) error {
	var decryptedKey bytes.Buffer
	iv, err := base64.StdEncoding.DecodeString(sealedKey.IV)
	if err != nil {
		return err
	}
	key, err := base64.StdEncoding.DecodeString(sealedKey.SealedKey)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, master)
	mac.Write(iv[:])
	mac.Write([]byte(sealedKey.BlobName))
	mac.Write([]byte("FROST-HMAC-SHA256"))

	config := sio.Config{MinVersion: sio.Version20, Key: mac.Sum(nil)}
	if n, err := sio.Decrypt(&decryptedKey, bytes.NewReader(key[:]), config); n != 32 || err != nil {
		return err
	}
	copy(k[:], decryptedKey.Bytes())
	return nil
}

func (k *Key) generateIV() (iv [32]byte) {
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		common.LogError("", err)
		panic(err)
	}
	return iv
}
