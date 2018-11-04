package config

type Config struct {
	Port        string
	Address     string
	PrivateKey  string
	PublicKey   string
	StoragePath string
}

var (
	Cfg = &Config{}
)

// func createHash(key string) string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(key))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func encrypt(data []byte, PublicKey string) []byte {
// 	block, _ := aes.NewCipher([]byte(createHash(PublicKey)))
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		panic(err.Error())
// 	}
// 	ciphertext := gcm.Seal(nonce, nonce, data, nil)
// 	return ciphertext
// }

// func decrypt(data []byte, PublicKey string) []byte {
// 	key := []byte(createHash(PublicKey))
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	nonceSize := gcm.NonceSize()
// 	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return plaintext
// }

// func encryptFile(filename string, data []byte, PublicKey string) {
// 	f, _ := os.Create(filename)
// 	defer f.Close()
// 	f.Write(encrypt(data, PublicKey))
// }

// func decryptFile(filename string, PublicKey string) []byte {
// 	data, _ := ioutil.ReadFile(filename)
// 	return decrypt(data, PublicKey)
// }
