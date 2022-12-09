package encrypt

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const AESKey = "asuperstrong32bitpasswordgohere!"

func GenHash(password string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Encrypt(hash string) (string, error) {
	key := []byte(AESKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	out := make([]byte, len(hash))
 
	block.Encrypt(out, []byte(hash))

	return hex.EncodeToString(out), nil
}

func Decrypt(cipherText string) (string, error) {
	text, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	key := []byte(AESKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(text))
	
    block.Decrypt(pt, text)
 
    s := string(pt[:])

	return s, nil
}