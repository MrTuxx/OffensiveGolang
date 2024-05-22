package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
)

// generateKeyFromPassword generates a 32-byte key from a given password using SHA-256
func generateKeyFromPassword(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// Random bytes for encryption key
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func handleEncryption(plaintext []byte, key []byte) (string, string) {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(key)
}

// Returns base64 encoded ciphertext encrypted with random 32-bit key
func GetEncryption(paylaod string) {

	data := []byte(paylaod)
	key_bytes, _ := generateRandomBytes(32)
	ciphertext, key := handleEncryption(data, key_bytes)

	println("[+] Encrypted paylaod: ", ciphertext)
	println("[+] Encrypted key: ", key)

	// the WriteFile method returns an error if unsuccessful
	ioutil.WriteFile("data.txt", []byte(ciphertext), 0777)
	println("[+] The file data.txt has been written with the base64 payload.")
}

// Returns base64 encoded ciphertext encrypted with 32-bit key
func GetEncryptionWithPassword(paylaod string, password string) {

	data := []byte(paylaod)
	key_bytes := generateKeyFromPassword(password)

	ciphertext, key := handleEncryption(data, key_bytes)
	println("[+] Encrypted paylaod: ", ciphertext)
	println("[+] Encrypted key: ", key)

	// the WriteFile method returns an error if unsuccessful
	ioutil.WriteFile("data.txt", []byte(ciphertext), 0777)
	println("[+] The file data.txt has been written with the base64 payload.")
}
