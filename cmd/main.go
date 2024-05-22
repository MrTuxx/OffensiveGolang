package main

import (
	"fmt"

	"github.com/MrTuxx/OffensiveGolang/pkg/encryption"
	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
)

func main() {

	//encryption.GetEncryption("<SHELLCODE>")

	url := "<URL>"
	pattern := `<REGEX TO GET THE PASSWORD IN THE URL>`
	password, err := evasion.ExtractMatchedStringFromURL(url, pattern)
	if err != nil {
		fmt.Println("[!] Error Extracting the String")
	}

	println("[+] Extracted Password: ", password)

	encryption.GetEncryptionWithPassword("<SHELLCODE>", password)

	println("[+] URL Brainfuck encoded: ", evasion.CodeBrainfuck(url))
	println("[+] Pattern Brainfuck encoded: ", evasion.CodeBrainfuck(pattern))

}
