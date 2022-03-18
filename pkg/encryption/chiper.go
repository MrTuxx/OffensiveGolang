package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
)

// Random bytes for encryption key
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Returns base64 encoded ciphertext encrypted with random 32-bit key
func encrypt(plaintext []byte) (string, string) {
	key, _ := generateRandomBytes(32)
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(key)
}

func GetEncryption(paylaod string) {

	data := []byte(paylaod)
	ciphertext, key := encrypt(data)
	println("[+] Encrypted paylaod: ", ciphertext)
	println("[+] Encrypted key: ", key)

	// the WriteFile method returns an error if unsuccessful
	ioutil.WriteFile("data.txt", []byte(ciphertext), 0777)
	println("[+] The file data.txt has been written with the base64 payload.")
}
