package util

import (
	"crypto/rand"
	"math/big"
	"strconv"
)

//Decrypt use to decrypt dest with rsa algoithm by private key
//Private key get from config
func Decrypt(dst interface{}) (interface{}, error) {
	return nil, nil
}

//Encrypt encrypt dst with public key
func Encrypt(dst interface{}) (interface{}, error) {
	return nil, nil
}

//EncryptWithKey encrypt dst with key
//Use AES algo or something like that
func EncryptWithKey(dst string, key interface{}) (string, error) {
	return "", nil
}

//DecryptWithKey decrypt dst with key
//Use AES algo or something like that
func DecryptWithKey(dst string, key interface{}) (string, error) {
	return "", nil
}

//GenerateSessionKey ...
func GenerateSessionKey() string {
	token := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		token = token + strconv.Itoa(int(n.Int64()))
	}
	return token
}
